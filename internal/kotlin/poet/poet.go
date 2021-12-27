package poet

import (
	"github.com/kyleconroy/sqlc/internal/kotlin/ast"
)

type proto interface {
	ProtoMessage()
}

func Node(node proto) *ast.Node {
	switch n := node.(type) {

	case *ast.Class:
		return &ast.Node{
			Node: &ast.Node_Class{
				Class: n,
			},
		}

	default:
		panic(n)

	}
}
