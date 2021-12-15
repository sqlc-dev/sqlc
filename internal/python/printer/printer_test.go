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
		"assign": {
			Node: &ast.Node{
				Node: &ast.Node_Assign{
					Assign: &ast.Assign{
						Targets: []*ast.Node{
							{
								Node: &ast.Node_Name{
									Name: &ast.Name{Id: "FICTION"},
								},
							},
						},
						Value: &ast.Node{
							Node: &ast.Node_Constant{
								Constant: &ast.Constant{
									Value: &ast.Constant_Str{
										Str: "FICTION",
									},
								},
							},
						},
					},
				},
			},
			Expected: `FICTION = "FICTION"`,
		},
		"class-base": {
			Node: &ast.Node{
				Node: &ast.Node_ClassDef{
					ClassDef: &ast.ClassDef{
						Name: "Foo",
						Bases: []*ast.Node{
							{
								Node: &ast.Node_Name{
									Name: &ast.Name{Id: "str"},
								},
							},
							{
								Node: &ast.Node_Attribute{
									Attribute: &ast.Attribute{
										Value: &ast.Node{
											Node: &ast.Node_Name{
												Name: &ast.Name{Id: "enum"},
											},
										},
										Attr: "Enum",
									},
								},
							},
						},
					},
				},
			},
			Expected: `class Foo(str, enum.Enum):`,
		},
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
													Slice: &ast.Node{
														Node: &ast.Node_Name{
															Name: &ast.Name{Id: "int"},
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
				},
			},
			Expected: `
@dataclass
class Foo:
    bar: int
    bat: Optional[int]
`,
		},
		"call": {
			Node: &ast.Node{
				Node: &ast.Node_Call{
					Call: &ast.Call{
						Func: &ast.Node{
							Node: &ast.Node_Alias{
								Alias: &ast.Alias{
									Name: "foo",
								},
							},
						},
					},
				},
			},
			Expected: `foo()`,
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
						Module: "pkg",
						Names: []*ast.Node{
							{
								Node: &ast.Node_Alias{
									Alias: &ast.Alias{
										Name: "foo",
									},
								},
							},
							{
								Node: &ast.Node_Alias{
									Alias: &ast.Alias{
										Name: "bar",
									},
								},
							},
						},
					},
				},
			},
			Expected: `from pkg import foo, bar`,
		},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			result := Print(tc.Node, Options{})
			if diff := cmp.Diff(strings.TrimSpace(tc.Expected), strings.TrimSpace(string(result.Python))); diff != "" {
				t.Errorf("print mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
