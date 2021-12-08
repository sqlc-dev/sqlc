package printer

import (
	"testing"

	"github.com/kyleconroy/sqlc/internal/python/ast"
)

func TestPrinter(t *testing.T) {
	node := &ast.Node{
		Node: &ast.Node_Module{
			Module: &ast.Module{
				Body: []*ast.Node{
					{
						Node: &ast.Node_Import{
							Import: &ast.Import{
								Names: []*ast.Node{
									{
										Node: &ast.Node_Alias{
											Alias: &ast.Alias{
												Name: "foo",
											},
										},
									},
								},
							},
						},
					},
					{
						Node: &ast.Node_ClassDef{
							ClassDef: &ast.ClassDef{Name: "Foo"},
						},
					},
				},
			},
		},
	}
	result := Print(node, Options{})
	t.Log(string(result.Code))
}
