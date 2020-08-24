package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterOpFamilyStmt struct {
	Opfamilyname *List
	Amname       *string
	IsDrop       bool
	Items        *List
}

func (n *AlterOpFamilyStmt) Pos() int {
	return 0
}
