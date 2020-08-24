package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DropStmt struct {
	Objects    *List
	RemoveType ObjectType
	Behavior   DropBehavior
	MissingOk  bool
	Concurrent bool
}

func (n *DropStmt) Pos() int {
	return 0
}
