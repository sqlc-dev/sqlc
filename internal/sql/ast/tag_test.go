package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// TestNoFieldShadowsTagKey guards the JSON encoding in json.go: every node
// object is emitted with a TagKey key naming its type, so no node may declare a
// field of its own that collides with it. encoding/json matches field names
// case-insensitively, so a field named "Tag" would silently capture the tag
// when the output is decoded back into a node.
//
// This is a unit test because the invariant is about the declarations in this
// package, not about anything sqlc produces. No SQL input can exercise a field
// that does not exist yet, so the end-to-end tests cannot catch the day someone
// adds one.
func TestNoFieldShadowsTagKey(t *testing.T) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, ".", nil, 0)
	if err != nil {
		t.Fatalf("parsing package: %s", err)
	}

	for _, pkg := range pkgs {
		ast.Inspect(pkg, func(n ast.Node) bool {
			spec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			structType, ok := spec.Type.(*ast.StructType)
			if !ok {
				return true
			}
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					if strings.EqualFold(name.Name, TagKey) {
						t.Errorf("%s declares a field %q, which collides with the %q key used to tag node types in JSON. Rename the field or pick a different TagKey.",
							spec.Name.Name, name.Name, TagKey)
					}
				}
			}
			return true
		})
	}
}
