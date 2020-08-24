package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreatedbStmt struct {
	Dbname  *string
	Options *List
}

func (n *CreatedbStmt) Pos() int {
	return 0
}
