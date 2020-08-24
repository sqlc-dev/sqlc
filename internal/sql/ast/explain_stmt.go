package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ExplainStmt struct {
	Query   ast.Node
	Options *List
}

func (n *ExplainStmt) Pos() int {
	return 0
}
