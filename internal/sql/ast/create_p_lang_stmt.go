package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreatePLangStmt struct {
	Replace     bool
	Plname      *string
	Plhandler   *List
	Plinline    *List
	Plvalidator *List
	Pltrusted   bool
}

func (n *CreatePLangStmt) Pos() int {
	return 0
}
