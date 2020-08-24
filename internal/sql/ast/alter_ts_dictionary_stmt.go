package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTSDictionaryStmt struct {
	Dictname *List
	Options  *List
}

func (n *AlterTSDictionaryStmt) Pos() int {
	return 0
}
