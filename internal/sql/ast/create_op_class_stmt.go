package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateOpClassStmt struct {
	Opclassname  *List
	Opfamilyname *List
	Amname       *string
	Datatype     *TypeName
	Items        *List
	IsDefault    bool
}

func (n *CreateOpClassStmt) Pos() int {
	return 0
}
