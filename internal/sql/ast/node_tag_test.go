package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// TestNodeTagFields guards the JSON encoding: every node struct must declare
// Tag NodeTag[T] as its first field with T naming the struct itself, and no
// node may declare another field that collides with the "tag" key. Nothing in
// the language enforces either — the field is copy-pasted between node files,
// so its type parameter can silently name the wrong node, and encoding/json
// matches field names case-insensitively on decode, so a second field named
// Tag would capture the tag string.
//
// This is a unit test because the invariant is about the declarations in this
// package, not about anything sqlc produces. No SQL input can exercise a node
// that has not been written yet, so the end-to-end tests cannot catch the day
// someone adds one without the field.
func TestNodeTagFields(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		t.Fatalf("parsing package: %s", err)
	}
	pkg := pkgs["ast"]

	// The node set: struct types with their own pointer-receiver Pos method.
	nodes := map[string]bool{}
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "Pos" || fn.Recv == nil {
				continue
			}
			if star, ok := fn.Recv.List[0].Type.(*ast.StarExpr); ok {
				if ident, ok := star.X.(*ast.Ident); ok {
					nodes[ident.Name] = true
				}
			}
		}
	}
	if len(nodes) == 0 {
		t.Fatal("found no Pos methods; the test is broken")
	}

	for _, file := range pkg.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			spec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			structType, ok := spec.Type.(*ast.StructType)
			if !ok || !nodes[spec.Name.Name] {
				return true
			}
			name := spec.Name.Name

			fields := structType.Fields.List
			if len(fields) == 0 || len(fields[0].Names) != 1 || fields[0].Names[0].Name != "Tag" {
				t.Errorf("%s: first field must be Tag NodeTag[%s]", name, name)
				return true
			}
			index, ok := fields[0].Type.(*ast.IndexExpr)
			if !ok {
				t.Errorf("%s: Tag field must have type NodeTag[%s]", name, name)
				return true
			}
			base, _ := index.X.(*ast.Ident)
			arg, _ := index.Index.(*ast.Ident)
			if base == nil || base.Name != "NodeTag" || arg == nil {
				t.Errorf("%s: Tag field must have type NodeTag[%s]", name, name)
				return true
			}
			if arg.Name != name {
				t.Errorf("%s: Tag field is NodeTag[%s]; its type parameter must name the containing struct", name, arg.Name)
			}

			for _, field := range fields[1:] {
				for _, fname := range field.Names {
					if strings.EqualFold(fname.Name, "tag") {
						t.Errorf("%s declares a second field %q, which collides with the tag field on JSON decode", name, fname.Name)
					}
				}
			}
			return true
		})
	}
}
