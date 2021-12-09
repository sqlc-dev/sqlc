package printer

import (
	"strings"
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
		"dataclass": {
			Node: &ast.Node{
				Node: &ast.Node_ClassDef{
					ClassDef: &ast.ClassDef{
						Name: "Foo",
						DecoratorList: []*ast.Node{
							{
								Node: &ast.Node_Name{
									Name: &ast.Name{
										Id: "dataclass",
									},
								},
							},
						},
						Body: []*ast.Node{
							{
								Node: &ast.Node_AnnAssign{
									AnnAssign: &ast.AnnAssign{
										Target: &ast.Name{Id: "bar"},
										Annotation: &ast.Node{
											Node: &ast.Node_Name{
												Name: &ast.Name{Id: "int"},
											},
										},
									},
								},
							},
							{
								Node: &ast.Node_AnnAssign{
									AnnAssign: &ast.AnnAssign{
										Target: &ast.Name{Id: "bat"},
										Annotation: &ast.Node{
											Node: &ast.Node_Subscript{
												Subscript: &ast.Subscript{
													Value: &ast.Name{Id: "Optional"},
													Slice: &ast.Name{Id: "int"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Expected: `
@dataclass
class Foo:
    bar: int
    bat: Optional[int]
`,
		},
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
			if diff := cmp.Diff(strings.TrimSpace(tc.Expected), strings.TrimSpace(string(result.Code))); diff != "" {
				t.Errorf("print mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
