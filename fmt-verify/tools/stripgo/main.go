// stripgo prints a Go file with every string-typed const/var initializer
// replaced by "", so generated code can be compared while ignoring the
// embedded SQL text (including `+ "`" +` concatenations).
package main

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
)

func isStringExpr(e ast.Expr) bool {
	switch v := e.(type) {
	case *ast.BasicLit:
		return v.Kind == token.STRING
	case *ast.BinaryExpr:
		return v.Op == token.ADD && isStringExpr(v.X) && isStringExpr(v.Y)
	case *ast.ParenExpr:
		return isStringExpr(v.X)
	}
	return false
}

func main() {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, os.Args[1], nil, 0)
	if err != nil {
		r, _ := os.Open(os.Args[1])
		if r != nil {
			io.Copy(os.Stdout, r)
		}
		return
	}
	ast.Inspect(f, func(n ast.Node) bool {
		gd, ok := n.(*ast.GenDecl)
		if !ok || (gd.Tok != token.CONST && gd.Tok != token.VAR) {
			return true
		}
		for _, spec := range gd.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			for i, v := range vs.Values {
				if isStringExpr(v) {
					vs.Values[i] = &ast.BasicLit{Kind: token.STRING, Value: `""`}
				}
			}
		}
		return true
	})
	printer.Fprint(os.Stdout, fset, f)
}
