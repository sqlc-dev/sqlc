package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterExtensionStmt struct {
	Extname *string
	Options *List
}

func (n *AlterExtensionStmt) Pos() int {
	return 0
}
