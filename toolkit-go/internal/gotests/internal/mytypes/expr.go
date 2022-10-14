package mytypes

import (
	"bytes"
	"fmt"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/gotests/internal/pkgMap"
	"go/ast"
	"path"
)

func ExprString(x ast.Expr, detail *pkgMap.PackageDetail) string {
	var buf bytes.Buffer
	WriteExpr(&buf, x, detail)
	return buf.String()
}

func WriteExpr(buf *bytes.Buffer, x ast.Expr, detail *pkgMap.PackageDetail) {

	switch x := x.(type) {
	default:
		buf.WriteString(fmt.Sprintf("(ast: %T)", x))

	case *ast.Ident:
		if pkgMap.CheckHasType(x, detail) {
			buf.WriteString(path.Base(detail.PkgName) + ".")
		}
		buf.WriteString(x.Name)

	case *ast.Ellipsis:
		buf.WriteString("...")
		if x.Elt != nil {
			WriteExpr(buf, x.Elt, detail)
		}

	case *ast.BasicLit:
		buf.WriteString(x.Value)

	case *ast.FuncLit:
		buf.WriteByte('(')
		WriteExpr(buf, x.Type, detail)
		buf.WriteString(" literal)")

	case *ast.CompositeLit:
		buf.WriteByte('(')
		WriteExpr(buf, x.Type, detail)
		buf.WriteString(" literal)")

	case *ast.ParenExpr:
		buf.WriteByte('(')
		WriteExpr(buf, x.X, detail)
		buf.WriteByte(')')

	case *ast.SelectorExpr:
		WriteExpr(buf, x.X, detail)
		buf.WriteByte('.')
		buf.WriteString(x.Sel.Name)

	case *ast.IndexExpr:
		WriteExpr(buf, x.X, detail)
		buf.WriteByte('[')
		exprs := UnpackExpr(x.Index)
		for i, e := range exprs {
			if i > 0 {
				buf.WriteString(", ")
			}
			WriteExpr(buf, e, detail)
		}
		buf.WriteByte(']')

	case *ast.SliceExpr:
		WriteExpr(buf, x.X, detail)
		buf.WriteByte('[')
		if x.Low != nil {
			WriteExpr(buf, x.Low, detail)
		}
		buf.WriteByte(':')
		if x.High != nil {
			WriteExpr(buf, x.High, detail)
		}
		if x.Slice3 {
			buf.WriteByte(':')
			if x.Max != nil {
				WriteExpr(buf, x.Max, detail)
			}
		}
		buf.WriteByte(']')

	case *ast.TypeAssertExpr:
		WriteExpr(buf, x.X, detail)
		buf.WriteString(".(")
		WriteExpr(buf, x.Type, detail)
		buf.WriteByte(')')

	case *ast.CallExpr:
		WriteExpr(buf, x.Fun, detail)
		buf.WriteByte('(')
		writeExprList(buf, x.Args, detail)
		if x.Ellipsis.IsValid() {
			buf.WriteString("...")
		}
		buf.WriteByte(')')

	case *ast.StarExpr:
		buf.WriteByte('*')
		WriteExpr(buf, x.X, detail)

	case *ast.UnaryExpr:
		buf.WriteString(x.Op.String())
		WriteExpr(buf, x.X, detail)

	case *ast.BinaryExpr:
		WriteExpr(buf, x.X, detail)
		buf.WriteByte(' ')
		buf.WriteString(x.Op.String())
		buf.WriteByte(' ')
		WriteExpr(buf, x.Y, detail)

	case *ast.ArrayType:
		buf.WriteByte('[')
		if x.Len != nil {
			WriteExpr(buf, x.Len, detail)
		}
		buf.WriteByte(']')
		WriteExpr(buf, x.Elt, detail)

	case *ast.StructType:
		buf.WriteString("struct{")
		writeFieldList(buf, x.Fields.List, "; ", false, detail)
		buf.WriteByte('}')

	case *ast.FuncType:
		buf.WriteString("func")
		writeSigExpr(buf, x, detail)

	case *ast.InterfaceType:
		var types []ast.Expr
		var methods []*ast.Field
		for _, f := range x.Methods.List {
			if len(f.Names) > 1 && f.Names[0].Name == "type" {
				types = append(types, f.Type)
			} else {
				methods = append(methods, f)
			}
		}

		buf.WriteString("interface{")
		writeFieldList(buf, methods, "; ", true, detail)
		if len(types) > 0 {
			if len(methods) > 0 {
				buf.WriteString("; ")
			}
			buf.WriteString("type ")
			writeExprList(buf, types, detail)
		}
		buf.WriteByte('}')

	case *ast.MapType:
		buf.WriteString("map[")
		WriteExpr(buf, x.Key, detail)
		buf.WriteByte(']')
		WriteExpr(buf, x.Value, detail)

	case *ast.ChanType:
		var s string
		switch x.Dir {
		case ast.SEND:
			s = "chan<- "
		case ast.RECV:
			s = "<-chan "
		default:
			s = "chan "
		}
		buf.WriteString(s)
		WriteExpr(buf, x.Value, detail)
	}
}

func UnpackExpr(expr ast.Expr) []ast.Expr {
	return []ast.Expr{expr}
}

func writeSigExpr(buf *bytes.Buffer, sig *ast.FuncType, detail *pkgMap.PackageDetail) {
	buf.WriteByte('(')
	writeFieldList(buf, sig.Params.List, ", ", false, detail)
	buf.WriteByte(')')

	res := sig.Results
	n := res.NumFields()
	if n == 0 {
		// no result
		return
	}

	buf.WriteByte(' ')
	if n == 1 && len(res.List[0].Names) == 0 {
		// single unnamed result
		WriteExpr(buf, res.List[0].Type, detail)
		return
	}

	// multiple or named result(s)
	buf.WriteByte('(')
	writeFieldList(buf, res.List, ", ", false, detail)
	buf.WriteByte(')')
}

func writeFieldList(buf *bytes.Buffer, list []*ast.Field, sep string, iface bool, detail *pkgMap.PackageDetail) {
	for i, f := range list {
		if i > 0 {
			buf.WriteString(sep)
		}

		writeIdentList(buf, f.Names)
		if sig, _ := f.Type.(*ast.FuncType); sig != nil && iface {
			writeSigExpr(buf, sig, detail)
			continue
		}
		if len(f.Names) > 0 {
			buf.WriteByte(' ')
		}

		WriteExpr(buf, f.Type, detail)

	}
}

func writeIdentList(buf *bytes.Buffer, list []*ast.Ident) {
	for i, x := range list {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(x.Name)
	}
}

func writeExprList(buf *bytes.Buffer, list []ast.Expr, detail *pkgMap.PackageDetail) {
	for i, x := range list {
		if i > 0 {
			buf.WriteString(", ")
		}
		WriteExpr(buf, x, detail)
	}
}
