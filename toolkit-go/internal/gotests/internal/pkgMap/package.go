package pkgMap

import (
	"go/ast"
	"go/token"
	"go/types"
)

type PackageDetail struct {
	Fs      []*ast.File
	Ti      map[ast.Expr]types.TypeAndValue
	PkgName string
}

func CheckHasType(expr ast.Expr, detail *PackageDetail) bool {
	if detail == nil {
		return false
	}
	for _, fAst := range detail.Fs {
		for _, d := range fAst.Decls {
			if decl, ok := d.(*ast.GenDecl); ok && decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					typeSpec, ok1 := spec.(*ast.TypeSpec)
					if !ok1 {
						continue
					}
					if types.ExprString(expr) == types.ExprString(typeSpec.Name) {
						return true
					}
				}

			}
			if decl, ok := d.(*ast.FuncDecl); ok {
				if decl.Recv == nil {
					if types.ExprString(expr) == types.ExprString(decl.Name) {
						return true
					}
				} else {
					// todo
				}
			}
		}
	}
	return false
}
