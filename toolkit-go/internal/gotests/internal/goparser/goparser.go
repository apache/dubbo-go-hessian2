// Package goparse contains logic for parsing Go files. Specifically it parses
// source and test files into domain models for generating tests.
package goparser

import (
	"errors"
	"fmt"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/internal/conf"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/internal/input"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/internal/mytypes"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/internal/pkgMap"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/util"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"path"
	"strings"

	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/internal/models"
)

// ErrEmptyFile represents an empty file error.
var ErrEmptyFile = errors.New("file is empty")

// Result representats a parsed Go file.
type Result struct {
	// The package name and imports of a Go file.
	Header *models.Header
	// All the functions and methods in a Go file.
	Funcs []*models.Function
}

// Parser can parse Go files.
type Parser struct {
	// The importer to resolve packages from import paths.
	Importer types.Importer
}

func (p *Parser) GetPackageParseMap(packageNames []string) (map[string]*pkgMap.PackageDetail, error) {
	goenv := util.New()
	res := make(map[string]*pkgMap.PackageDetail, 0)
	fset := token.NewFileSet()
	conf := &types.Config{
		Importer: p.Importer,
		// Adding a NO-OP error function ignores errors and performs best-effort
		// type checking. https://godoc.org/golang.org/x/tools/go/types#Config
		Error: func(error) {},
	}
	for _, packageName := range packageNames {
		dir, err := goenv.FindPackageDir(packageName)
		if err == nil && len(dir) > 0 {
			files, err := input.Files(dir)
			if err == nil && len(files) > 0 {
				fs, err := p.parseFilesByFlies(fset, files)
				if err == nil {
					ti := &types.Info{
						Types: make(map[ast.Expr]types.TypeAndValue),
					}
					// Note: conf.Check can fail, but since Info is not required data, it's ok.
					conf.Check("", fset, fs, ti)
					res[packageName] = &pkgMap.PackageDetail{
						Fs:      fs,
						Ti:      ti.Types,
						PkgName: packageName,
					}
				}
			}
		}
	}
	return res, nil
}

func (p *Parser) parseFilesByFlies(fset *token.FileSet, files []models.Path) ([]*ast.File, error) {
	var fs []*ast.File
	for _, file := range files {
		ff, err := parser.ParseFile(fset, string(file), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("other file parser.ParseFile: %v", err)
		}
		fs = append(fs, ff)
	}
	return fs, nil
}

// Parse parses a given Go file at srcPath, along any files that share the same
// pkgMap, into a domain model for generating tests.
func (p *Parser) Parse(srcPath string, files []models.Path) (*Result, error) {
	b, err := p.readFile(srcPath)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	f, err := p.parseFile(fset, srcPath)
	if err != nil {
		return nil, err
	}
	fs, err := p.parseFiles(fset, f, files)
	if err != nil {
		return nil, err
	}
	imports := parseImports(f.Imports)
	packageAstMap := make(map[string]*pkgMap.PackageDetail)
	if len(imports) > 0 {
		packages := make([]string, 0, len(imports))
		for _, item := range imports {
			packages = append(packages, strings.Trim(item.Path, "\""))
		}
		packageAstMap, _ = p.GetPackageParseMap(packages)
	}
	return &Result{
		Header: &models.Header{
			Comments: parsePkgComment(f, f.Package),
			Package:  f.Name.String(),
			Imports:  imports,
			Code:     goCode(b, f),
		},
		Funcs: p.parseFunctions(fset, f, fs, packageAstMap),
	}, nil
}

func (p *Parser) readFile(srcPath string) ([]byte, error) {
	b, err := ioutil.ReadFile(srcPath)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadFile: %v", err)
	}
	if len(b) == 0 {
		return nil, ErrEmptyFile
	}
	return b, nil
}

func (p *Parser) parseFile(fset *token.FileSet, srcPath string) (*ast.File, error) {
	f, err := parser.ParseFile(fset, srcPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("target parser.ParseFile(): %v", err)
	}
	return f, nil
}

func (p *Parser) parseFiles(fset *token.FileSet, f *ast.File, files []models.Path) ([]*ast.File, error) {
	pkg := f.Name.String()
	var fs []*ast.File
	for _, file := range files {
		ff, err := parser.ParseFile(fset, string(file), nil, 0)
		if err != nil {
			return nil, fmt.Errorf("other file parser.ParseFile: %v", err)
		}
		if name := ff.Name.String(); name != pkg {
			continue
		}
		fs = append(fs, ff)
	}
	return fs, nil
}

func (p *Parser) parseFunctions(fset *token.FileSet, f *ast.File, fs []*ast.File, packageAstMap map[string]*pkgMap.PackageDetail) []*models.Function {
	ul, el, exprTypeMap := p.parseTypes(fset, fs)
	//for _, pkgAst := range packageAstMap {
	//	for e, t := range pkgAst.Ti {
	//		ul[t.Type.String()] = t.Type.Underlying()
	//		if v, ok := t.Type.(*types.Struct); ok {
	//			el[v] = e
	//		}
	//	}
	//}

	importers := parseImports(f.Imports)
	var funcs []*models.Function
	for _, d := range f.Decls {
		fDecl, ok := d.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcs = append(funcs, parseFunc(fDecl, ul, el, exprTypeMap, importers, packageAstMap))
	}
	return funcs
}

func (p *Parser) parseTypes(fset *token.FileSet, fs []*ast.File) (map[string]types.Type, map[*types.Struct]ast.Expr, map[ast.Expr]types.TypeAndValue) {
	conf := &types.Config{
		Importer: p.Importer,
		// Adding a NO-OP error function ignores errors and performs best-effort
		// type checking. https://godoc.org/golang.org/x/tools/go/types#Config
		Error: func(error) {},
	}
	ti := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}
	// Note: conf.Check can fail, but since Info is not required data, it's ok.
	conf.Check("", fset, fs, ti)
	ul := make(map[string]types.Type)
	el := make(map[*types.Struct]ast.Expr)
	typeStrList := []string{}
	for e, t := range ti.Types {
		// Collect the underlying mytypes.
		ul[t.Type.String()] = t.Type.Underlying()
		//if res, ok := e.(*ast.CallExpr); ok {
		//	fmt.Println(e, "\t", ti.Types[res.Fun].Type.String())
		//}

		typeStrList = append(typeStrList, t.Type.String())
		// Collect structs to determine the fields of a receiver.
		if v, ok := t.Type.(*types.Struct); ok {
			el[v] = e
		}
	}
	return ul, el, ti.Types
}

func parsePkgComment(f *ast.File, pkgPos token.Pos) []string {
	var comments []string
	var count int

	for _, comment := range f.Comments {

		if comment.End() >= pkgPos {
			break
		}
		for _, c := range comment.List {
			count += len(c.Text) + 1 // +1 for '\n'
			if count < int(c.End()) {
				n := int(c.End()) - count - 1
				comments = append(comments, strings.Repeat("\n", n))
				count++ // for last of '\n'
			}
			comments = append(comments, c.Text)
		}
	}

	if int(pkgPos)-count > 1 {
		comments = append(comments, strings.Repeat("\n", int(pkgPos)-count-2))
	}
	return comments
}

// Returns the Go code below the imports block.
func goCode(b []byte, f *ast.File) []byte {
	furthestPos := f.Name.End()
	for _, node := range f.Imports {
		if pos := node.End(); pos > furthestPos {
			furthestPos = pos
		}
	}
	if furthestPos < token.Pos(len(b)) {
		furthestPos++

		// Avoid wrong output on windows-encoded files
		if b[furthestPos-2] == '\r' && b[furthestPos-1] == '\n' && furthestPos < token.Pos(len(b)) {
			furthestPos++
		}
	}
	return b[furthestPos:]
}

func parseFunc(
	fDecl *ast.FuncDecl,
	ul map[string]types.Type,
	el map[*types.Struct]ast.Expr,
	exprTypeMap map[ast.Expr]types.TypeAndValue,
	imports []*models.Import,
	packageAstMap map[string]*pkgMap.PackageDetail) *models.Function {
	f := &models.Function{
		Name:       fDecl.Name.String(),
		IsExported: fDecl.Name.IsExported(),
		Receiver:   parseReceiver(fDecl.Recv, ul, el, packageAstMap),
		Parameters: parseFieldList(fDecl.Type.Params, ul, packageAstMap),
	}
	fs := parseFieldList(fDecl.Type.Results, ul, packageAstMap)
	i := 0
	for _, fi := range fs {
		if fi.Type.String() == "error" {
			f.ReturnsError = true
			continue
		}
		fi.Index = i
		f.Results = append(f.Results, fi)
		i++
	}
	if fDecl.Body != nil && len(fDecl.Body.List) > 0 {
		f.MockFuncList = parseMockFunc(fDecl.Body.List, exprTypeMap, imports, packageAstMap)
	}
	return f
}

func GetFunExprString(expr ast.Expr) (res string) {
	switch expr.(type) {
	case *ast.SelectorExpr:
		funExpr := expr.(*ast.SelectorExpr)
		if funX, ok := funExpr.X.(*ast.Ident); ok {
			res += funX.String()
		}
		res += "." + funExpr.Sel.String()
	case *ast.Ident:
		funExpr := expr.(*ast.Ident)
		res += funExpr.String()
	}
	return
}

func ReplacePackage(str string, imports []*models.Import) string {
	for _, item := range imports {
		path := strings.Trim(item.Path, "\"")
		if strings.Contains(str, path+".") {
			name := item.Name
			if name == "" {
				paths := strings.Split(path, "/")
				if len(paths) > 0 {
					name = paths[len(paths)-1]
				}
			}
			str = strings.Replace(str, path+".", name+".", -1)
		}
	}
	return str
}

func parseMockFunc(bodyList []ast.Stmt, exprTypeMap map[ast.Expr]types.TypeAndValue, imports []*models.Import, packageAstMap map[string]*pkgMap.PackageDetail) (res []*models.MockFuncItem) {
	exprMap := make(map[string]string)
	for expr := range exprTypeMap {
		exprMap[types.ExprString(expr)] = exprTypeMap[expr].Type.String()
	}

	for _, stmt := range bodyList {
		res = append(res, addStmtMockFunc(stmt, exprMap, packageAstMap)...)
	}
	for _, mockStr := range res {
		mockStr.MockFuncStr = ReplacePackage(mockStr.MockFuncStr, imports)
	}
	return res
}

func addStmtMockFunc(stmt ast.Stmt, exprTypeMap map[string]string, packageAstMap map[string]*pkgMap.PackageDetail) (res []*models.MockFuncItem) {
	if stmt == nil {
		return
	}
	switch stmt.(type) {
	case *ast.AssignStmt:
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			return
		}
		for _, e := range assignStmt.Lhs {
			res = append(res, addExprMockFunc(e, exprTypeMap, packageAstMap)...)
		}
		for _, e := range assignStmt.Rhs {
			res = append(res, addExprMockFunc(e, exprTypeMap, packageAstMap)...)
		}
	case *ast.IfStmt:
		ifStmt, ok := stmt.(*ast.IfStmt)
		if !ok {
			return
		}

		res = append(res, addExprMockFunc(ifStmt.Cond, exprTypeMap, packageAstMap)...)
		res = append(res, addStmtMockFunc(ifStmt.Init, exprTypeMap, packageAstMap)...)
		res = append(res, addStmtMockFunc(ifStmt.Else, exprTypeMap, packageAstMap)...)
		if ifStmt.Body != nil {
			for _, s := range ifStmt.Body.List {
				res = append(res, addStmtMockFunc(s, exprTypeMap, packageAstMap)...)
			}
		}
	case *ast.ReturnStmt:
		returnStmt, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			return
		}
		for _, e := range returnStmt.Results {
			res = append(res, addExprMockFunc(e, exprTypeMap, packageAstMap)...)
		}
	case *ast.ForStmt:
		forStmt, ok := stmt.(*ast.ForStmt)
		if !ok {
			return
		}
		res = append(res, addExprMockFunc(forStmt.Cond, exprTypeMap, packageAstMap)...)
		res = append(res, addStmtMockFunc(forStmt.Init, exprTypeMap, packageAstMap)...)
		res = append(res, addStmtMockFunc(forStmt.Post, exprTypeMap, packageAstMap)...)
		if forStmt.Body != nil {
			for _, s := range forStmt.Body.List {
				res = append(res, addStmtMockFunc(s, exprTypeMap, packageAstMap)...)
			}
		}
	case *ast.ExprStmt:
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			return
		}
		res = append(res, addExprMockFunc(exprStmt.X, exprTypeMap, packageAstMap)...)
	default:
		//todo other Stmt
	}
	return
}

func GetFunExprStringAndType(expr ast.Expr) (res, xName string, X *ast.Ident) {
	switch expr.(type) {
	case *ast.SelectorExpr:
		funExpr := expr.(*ast.SelectorExpr)
		if funX, ok := funExpr.X.(*ast.Ident); ok {
			if funX.Obj != nil {
				X = funX
				res += funExpr.Sel.String()
			} else {
				res += funX.String()
				xName = funX.String()
				res += "." + funExpr.Sel.String()
			}
		}
	case *ast.Ident:
		funExpr := expr.(*ast.Ident)
		res += funExpr.String()
	}
	return
}

func addExprMockFunc(expr ast.Expr, exprTypeMap map[string]string, packageAstMap map[string]*pkgMap.PackageDetail) (res []*models.MockFuncItem) {
	if expr == nil {
		return
	}
	switch expr.(type) {
	case *ast.CallExpr:
		callExp, ok := expr.(*ast.CallExpr)
		if !ok {
			return
		}
		if _, exist := exprTypeMap[GetFunExprString(callExp.Fun)]; exist {
			exprStr, xName, X := GetFunExprStringAndType(callExp.Fun)
			if !conf.CheckInPkg(xName) {
				if exprStr != "append" && exprStr != "make" && exprStr != "len" {
					if X == nil {
						res = append(res, &models.MockFuncItem{
							MockFuncStr: fmt.Sprintf("//monkey.Patch(%s,%s{\n\t//\tpanic(\"should impl\")\n\t//})", types.ExprString(callExp.Fun), exprTypeMap[types.ExprString(callExp.Fun)]),
						})
					} else {
						objStr, ok := exprTypeMap[X.String()]
						if ok {
							res = append(res, &models.MockFuncItem{
								MockFuncStr: fmt.Sprintf("//var %s %s\n//monkey.PatchInstanceMethod(reflect.TypeOf(%s),\"%s\",%s{\n\t//\tpanic(\"should impl\")\n\t//})", X.String(), objStr, X.String(), exprStr, exprTypeMap[types.ExprString(callExp.Fun)]),
							})
						}
					}

				}
			}
		} else {
			if funExpr, ok := callExp.Fun.(*ast.SelectorExpr); ok {
				if funX, ok := funExpr.X.(*ast.Ident); ok {
					for _, pkgFileAst := range packageAstMap {
						for _, fileAst := range pkgFileAst.Fs {
							if funX.String() != fileAst.Name.Name {
								break
							}
							if !conf.CheckInPkg(funX.String()) {
								break
							}
							for _, decl := range fileAst.Decls {
								if funcDecl, ok := decl.(*ast.FuncDecl); ok && funExpr.Sel.String() == funcDecl.Name.Name {
									str := "func "
									if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
										str += "("
										paramsList := funcDecl.Type.Params.List
										for i, params := range paramsList {
											if i > 0 {
												str += ", "
											}
											for j, name := range params.Names {
												if j > 0 {
													str += ", "
												}
												str += mytypes.ExprString(name, pkgFileAst) + " "
											}
											str += mytypes.ExprString(params.Type, pkgFileAst)
										}
										str += ") "
									} else {
										str += "() "
									}
									if funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0 {
										resultsList := funcDecl.Type.Results.List
										if len(resultsList) > 1 {
											str += "("
										}
										for i, results := range resultsList {
											if i > 0 {
												str += ", "
											}
											str += mytypes.ExprString(results.Type, pkgFileAst)
										}
										if len(resultsList) > 1 {
											str += ") "
										}
									}
									res = append(res, &models.MockFuncItem{
										MockFuncStr: fmt.Sprintf("//monkey.Patch(%s,%s{\n\t//\tpanic(\"should impl\")\n\t//})", types.ExprString(callExp.Fun), str),
									})
									break
								}
							}
						}
					}
				}
			}

		}
		for _, e := range callExp.Args {
			res = append(res, addExprMockFunc(e, exprTypeMap, packageAstMap)...)
		}
	default:
		//todo other Expr
	}
	return
}

func parseImports(imps []*ast.ImportSpec) []*models.Import {
	var is []*models.Import
	for _, imp := range imps {
		var n string
		if imp.Name != nil {
			n = imp.Name.String()
		}
		is = append(is, &models.Import{
			Name: n,
			Path: imp.Path.Value,
		})
	}
	return is
}

func parseReceiver(fl *ast.FieldList, ul map[string]types.Type, el map[*types.Struct]ast.Expr, packageAstMap map[string]*pkgMap.PackageDetail) *models.Receiver {
	if fl == nil {
		return nil
	}
	r := &models.Receiver{
		Field: parseFieldList(fl, ul, packageAstMap)[0],
	}
	t, ok := ul[r.Type.Value]
	if !ok {
		return r
	}
	s, ok := t.(*types.Struct)
	if !ok {
		return r
	}
	st, found := el[s]
	if !found {
		return r
	}
	r.Fields = append(r.Fields, parseFieldList(st.(*ast.StructType).Fields, ul, packageAstMap)...)
	for i, f := range r.Fields {
		if i >= s.NumFields() {
			break
		}
		f.Name = s.Field(i).Name()
	}
	return r

}

func findObj(fType ast.Expr, packageAstMap map[string]*pkgMap.PackageDetail, identSelfPkgMap map[string]string) {
	switch fType.(type) {
	case *ast.SelectorExpr:
		if expr, ok := fType.(*ast.SelectorExpr); ok {
			if x, ok1 := expr.X.(*ast.Ident); ok1 && expr.Sel.Obj == nil {
				pkgName := x.String()
				for pkg, packageAst := range packageAstMap {
					if path.Base(pkg) == pkgName {
						for _, fAst := range packageAst.Fs {
							for _, d := range fAst.Decls {
								if decl, ok2 := d.(*ast.GenDecl); ok2 {

									for _, s := range decl.Specs {
										if typeSpec, ok3 := s.(*ast.TypeSpec); ok3 && typeSpec.Name.Name == expr.Sel.Name {
											expr.Sel.Obj = typeSpec.Name.Obj
											identSelfPkgMap[expr.Sel.Name] = path.Base(pkg)
											return
										}
									}
								}
							}
						}
					}
				}
			}
		}
	case *ast.StarExpr:
		if expr, ok := fType.(*ast.StarExpr); ok {
			findObj(expr.X, packageAstMap, identSelfPkgMap)
		}
	case *ast.ArrayType:
		if expr, ok := fType.(*ast.ArrayType); ok {
			findObj(expr.Elt, packageAstMap, identSelfPkgMap)
		}
	case *ast.MapType:
		if expr, ok := fType.(*ast.MapType); ok {
			findObj(expr.Key, packageAstMap, identSelfPkgMap)
			findObj(expr.Value, packageAstMap, identSelfPkgMap)
		}
	}
}

func parseFieldList(fl *ast.FieldList, ul map[string]types.Type, packageAstMap map[string]*pkgMap.PackageDetail) []*models.Field {
	if fl == nil {
		return nil
	}
	i := 0
	var fs []*models.Field
	for _, f := range fl.List {
		identSelfPkgMap := make(map[string]string)
		findObj(f.Type, packageAstMap, identSelfPkgMap)

		for _, pf := range parseFields(f, ul, packageAstMap, identSelfPkgMap) {
			pf.Index = i
			fs = append(fs, pf)
			i++
		}
	}
	return fs
}

func parseFields(f *ast.Field, ul map[string]types.Type, packageAstMap map[string]*pkgMap.PackageDetail, identSelfPkgMap map[string]string) []*models.Field {
	t := parseExpr(f.Type, ul, packageAstMap, identSelfPkgMap)
	if len(f.Names) == 0 {
		return []*models.Field{{
			Type: t,
		}}
	}
	var fs []*models.Field
	for _, n := range f.Names {
		fs = append(fs, &models.Field{
			Name: n.Name,
			Type: t,
		})
	}
	return fs
}

func parseExpr(e ast.Expr, ul map[string]types.Type, packageAstMap map[string]*pkgMap.PackageDetail, identSelfPkgMap map[string]string) *models.Expression {
	switch v := e.(type) {
	case *ast.StarExpr:
		val := types.ExprString(v.X)
		//exp := parseExpr(v.X, ul, packageAstMap, identSelfPkgMap)
		return &models.Expression{
			Value:      val,
			IsStar:     true,
			Underlying: underlying(val, ul),
			//ChildField:      exp.ChildField,
			IdentSelfPkgMap: identSelfPkgMap,
		}
	case *ast.Ellipsis:
		exp := parseExpr(v.Elt, ul, packageAstMap, identSelfPkgMap)
		return &models.Expression{
			Value:      exp.Value,
			IsStar:     exp.IsStar,
			IsVariadic: true,
			Underlying: underlying(exp.Value, ul),
			//ChildField:      exp.ChildField,
			IdentSelfPkgMap: identSelfPkgMap,
		}
	case *ast.InterfaceType:
		val := types.ExprString(e)
		return &models.Expression{
			Value:           val,
			Underlying:      underlying(val, ul),
			IsWriter:        val == "io.Writer",
			IsInterface:     true,
			IdentSelfPkgMap: identSelfPkgMap,
		}
	case *ast.StructType:
		//expr := e.(*ast.StructType)
		val := types.ExprString(e)
		//childFields := parseFieldList(expr.Fields, ul, packageAstMap)

		return &models.Expression{
			Value:       val,
			Underlying:  underlying(val, ul),
			IsWriter:    val == "io.Writer",
			IsInterface: false,
			//ChildField:      childFields,
			IdentSelfPkgMap: identSelfPkgMap,
		}
	case *ast.Ident:
		isInterface := false
		val := types.ExprString(e)
		exprType, ok := ul[val]
		if ok {
			_, isInterface = exprType.(*types.Interface)
		}
		//var childFields []*models.Field
		//expr := e.(*ast.Ident)
		//if expr.Obj != nil {
		//	decl := expr.Obj.Decl.(*ast.TypeSpec)
		//	exp := parseExpr(decl.Type, ul, packageAstMap, identSelfPkgMap)
		//	childFields = exp.ChildField
		//}
		funStr := ""
		isFun := false
		expr := e.(*ast.Ident)
		if expr.Obj != nil {
			if typeSpec, ok := expr.Obj.Decl.(*ast.TypeSpec); ok {
				switch v := typeSpec.Type.(type) {
				case *ast.FuncType:
					isFun = true
					str := "func "
					if v.Params != nil && len(v.Params.List) > 0 {
						str += "("
						paramsList := v.Params.List
						for i, params := range paramsList {
							if i > 0 {
								str += ", "
							}
							for j, name := range params.Names {
								if j > 0 {
									str += ", "
								}
								str += mytypes.ExprString(name, nil) + " "
							}
							str += mytypes.ExprString(params.Type, nil)
						}
						str += ") "
					} else {
						str += "() "
					}
					if v.Results != nil && len(v.Results.List) > 0 {
						resultsList := v.Results.List
						if len(resultsList) > 1 {
							str += "("
						}
						for i, results := range resultsList {
							if i > 0 {
								str += ", "
							}
							str += mytypes.ExprString(results.Type, nil)
						}
						if len(resultsList) > 1 {
							str += ") "
						}
					}
					funStr = str
				case *ast.InterfaceType:
					isInterface = true
				}
			}
		}

		return &models.Expression{
			Value:       val,
			Underlying:  underlying(val, ul),
			IsWriter:    val == "io.Writer",
			IsInterface: isInterface,
			//ChildField:      childFields,
			IdentSelfPkgMap: identSelfPkgMap,
			IsFuncType:      isFun,
			FuncStr:         funStr,
		}
	case *ast.SelectorExpr:
		isInterface := false
		val := types.ExprString(e)
		exprType, ok := ul[val]
		if ok {
			_, isInterface = exprType.(*types.Interface)
		}
		//var childFields []*models.Field
		//expr := e.(*ast.SelectorExpr)
		//if expr.Sel.Obj != nil {
		//	decl := expr.Sel.Obj.Decl.(*ast.TypeSpec)
		//	exp := parseExpr(decl.Type, ul, packageAstMap, identSelfPkgMap)
		//	childFields = exp.ChildField
		//}
		funStr := ""
		isFun := false
		expr := v.Sel
		pkgName := ""
		if xInfo, ok := v.X.(*ast.Ident); ok {
			pkgName = xInfo.Name
		}
		if expr.Obj != nil {
			if typeSpec, ok := expr.Obj.Decl.(*ast.TypeSpec); ok {
				switch v := typeSpec.Type.(type) {
				case *ast.FuncType:
					var detail *pkgMap.PackageDetail
					for pkg, packageAst := range packageAstMap {
						if path.Base(pkg) == pkgName {
							detail = packageAst
							break
						}
					}
					isFun = true
					str := "func "
					if v.Params != nil && len(v.Params.List) > 0 {
						str += "("
						paramsList := v.Params.List
						for i, params := range paramsList {
							if i > 0 {
								str += ", "
							}
							for j, name := range params.Names {
								if j > 0 {
									str += ", "
								}
								str += mytypes.ExprString(name, detail) + " "
							}
							str += mytypes.ExprString(params.Type, detail)
						}
						str += ") "
					} else {
						str += "() "
					}
					if v.Results != nil && len(v.Results.List) > 0 {
						resultsList := v.Results.List
						if len(resultsList) > 1 {
							str += "("
						}
						for i, results := range resultsList {
							if i > 0 {
								str += ", "
							}
							str += mytypes.ExprString(results.Type, detail)
						}
						if len(resultsList) > 1 {
							str += ") "
						}
					}
					funStr = str
				case *ast.InterfaceType:
					isInterface = true
				}
			}
		}
		return &models.Expression{
			Value:       val,
			Underlying:  underlying(val, ul),
			IsWriter:    val == "io.Writer",
			IsInterface: isInterface,
			//ChildField:      childFields,
			IdentSelfPkgMap: identSelfPkgMap,
			IsFuncType:      isFun,
			FuncStr:         funStr,
		}
	case *ast.FuncType:
		str := "func "
		if v.Params != nil && len(v.Params.List) > 0 {
			str += "("
			paramsList := v.Params.List
			for i, params := range paramsList {
				if i > 0 {
					str += ", "
				}
				for j, name := range params.Names {
					if j > 0 {
						str += ", "
					}
					str += mytypes.ExprString(name, nil) + " "
				}
				str += mytypes.ExprString(params.Type, nil)
			}
			str += ") "
		} else {
			str += "() "
		}
		if v.Results != nil && len(v.Results.List) > 0 {
			resultsList := v.Results.List
			if len(resultsList) > 1 {
				str += "("
			}
			for i, results := range resultsList {
				if i > 0 {
					str += ", "
				}
				str += mytypes.ExprString(results.Type, nil)
			}
			if len(resultsList) > 1 {
				str += ") "
			}
		}
		isInterface := false
		val := types.ExprString(e)
		exprType, ok := ul[val]
		if ok {
			_, isInterface = exprType.(*types.Interface)
		}
		return &models.Expression{
			Value:           val,
			Underlying:      underlying(val, ul),
			IsWriter:        val == "io.Writer",
			IsInterface:     isInterface,
			IdentSelfPkgMap: identSelfPkgMap,
			IsFuncType:      true,
			FuncStr:         str,
		}
	default:
		isInterface := false
		val := types.ExprString(e)
		exprType, ok := ul[val]
		if ok {
			_, isInterface = exprType.(*types.Interface)
		}
		return &models.Expression{
			Value:           val,
			Underlying:      underlying(val, ul),
			IsWriter:        val == "io.Writer",
			IsInterface:     isInterface,
			IdentSelfPkgMap: identSelfPkgMap,
		}
	}
}

func underlying(val string, ul map[string]types.Type) string {
	if ul[val] != nil {
		return ul[val].String()
	}
	return ""
}
