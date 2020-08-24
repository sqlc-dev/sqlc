package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ConstraintsSetStmt struct {
	Constraints *List
	Deferred    bool
}

func (n *ConstraintsSetStmt) Pos() int {
	return 0
}
