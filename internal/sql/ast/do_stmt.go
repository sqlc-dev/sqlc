package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DoStmt struct {
	Args *List
}

func (n *DoStmt) Pos() int {
	return 0
}
