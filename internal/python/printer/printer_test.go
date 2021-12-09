package printer

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kyleconroy/sqlc/internal/python/ast"
)

type testcase struct {
	Node     *ast.Node
	Expected string
}

func TestPrinter(t *testing.T) {
	for name, tc := range map[string]testcase{
		"import": {
			Node: &ast.Node{
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
			Expected: `import foo`,
		},
		"import-from": {
			Node: &ast.Node{
				Node: &ast.Node_ImportFrom{
					ImportFrom: &ast.ImportFrom{
						Module: "dataclasses",
						Names: []*ast.Node{
							{
								Node: &ast.Node_Alias{
									Alias: &ast.Alias{
										Name: "dataclass",
									},
								},
							},
						},
					},
				},
			},
			Expected: `from dataclasses import dataclass`,
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result := Print(tc.Node, Options{})
			if diff := cmp.Diff(tc.Expected, string(result.Code)); diff != "" {
				t.Errorf("print mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
