package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTableStmt struct {
	Relation  *RangeVar
	Cmds      *List
	Relkind   ObjectType
	MissingOk bool
}

func (n *AlterTableStmt) Pos() int {
	return 0
}
