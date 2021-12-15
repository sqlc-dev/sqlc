package poet

import (
	"github.com/kyleconroy/sqlc/internal/python/ast"
)

type proto interface {
	ProtoMessage()
}

func Nodes(nodes ...proto) []*ast.Node {
	list := make([]*ast.Node, len(nodes))
	for i, _ := range nodes {
		list[i] = Node(nodes[i])
	}
	return list
}

func Node(node proto) *ast.Node {
	switch n := node.(type) {

	case *ast.Alias:
		return &ast.Node{
			Node: &ast.Node_Alias{
				Alias: n,
			},
		}

	case *ast.Await:
		return &ast.Node{
			Node: &ast.Node_Await{
				Await: n,
			},
		}

	case *ast.AnnAssign:
		return &ast.Node{
			Node: &ast.Node_AnnAssign{
				AnnAssign: n,
			},
		}

	case *ast.Assign:
		return &ast.Node{
			Node: &ast.Node_Assign{
				Assign: n,
			},
		}

	case *ast.AsyncFor:
		return &ast.Node{
			Node: &ast.Node_AsyncFor{
				AsyncFor: n,
			},
		}

	case *ast.AsyncFunctionDef:
		return &ast.Node{
			Node: &ast.Node_AsyncFunctionDef{
				AsyncFunctionDef: n,
			},
		}

	case *ast.Attribute:
		return &ast.Node{
			Node: &ast.Node_Attribute{
				Attribute: n,
			},
		}

	case *ast.Call:
		return &ast.Node{
			Node: &ast.Node_Call{
				Call: n,
			},
		}

	case *ast.ClassDef:
		return &ast.Node{
			Node: &ast.Node_ClassDef{
				ClassDef: n,
			},
		}

	case *ast.Comment:
		return &ast.Node{
			Node: &ast.Node_Comment{
				Comment: n,
			},
		}

	case *ast.Compare:
		return &ast.Node{
			Node: &ast.Node_Compare{
				Compare: n,
			},
		}

	// case *ast.Constant:

	// case *ast.Dict:

	case *ast.Expr:
		return &ast.Node{
			Node: &ast.Node_Expr{
				Expr: n,
			},
		}

	case *ast.For:
		return &ast.Node{
			Node: &ast.Node_For{
				For: n,
			},
		}

	case *ast.FunctionDef:
		return &ast.Node{
			Node: &ast.Node_FunctionDef{
				FunctionDef: n,
			},
		}

	case *ast.If:
		return &ast.Node{
			Node: &ast.Node_If{
				If: n,
			},
		}

	// case *ast.Node_Import:
	// 	w.printImport(n.Import, indent)

	// case *ast.Node_ImportFrom:
	// 	w.printImportFrom(n.ImportFrom, indent)

	// case *ast.Node_Is:
	// 	w.print("is")

	// case *ast.Node_Keyword:
	// 	w.printKeyword(n.Keyword, indent)

	case *ast.Module:
		return &ast.Node{
			Node: &ast.Node_Module{
				Module: n,
			},
		}

	// 	w.printModule(n.Module, indent)

	// case *ast.Node_Name:
	// 	w.print(n.Name.Id)

	// case *ast.Node_Pass:
	// 	w.print("pass")

	// case *ast.Node_Return:
	// 	w.printReturn(n.Return, indent)

	// case *ast.Node_Subscript:
	// 	w.printSubscript(n.Subscript, indent)

	case *ast.Yield:
		return &ast.Node{
			Node: &ast.Node_Yield{
				Yield: n,
			},
		}

	default:
		panic(n)
	}

}

func Constant(value interface{}) *ast.Node {
	switch n := value.(type) {
	case string:
		return &ast.Node{
			Node: &ast.Node_Constant{
				Constant: &ast.Constant{
					Value: &ast.Constant_Str{
						Str: n,
					},
				},
			},
		}

	case int:
		return &ast.Node{
			Node: &ast.Node_Constant{
				Constant: &ast.Constant{
					Value: &ast.Constant_Int{
						Int: int32(n),
					},
				},
			},
		}

	case nil:
		return &ast.Node{
			Node: &ast.Node_Constant{
				Constant: &ast.Constant{
					Value: &ast.Constant_None{},
				},
			},
		}

	default:
		panic("unknown type")
	}
}
