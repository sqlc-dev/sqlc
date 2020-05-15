package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateOpClassStmt struct {
	Opclassname  *ast.List
	Opfamilyname *ast.List
	Amname       *string
	Datatype     *TypeName
	Items        *ast.List
	IsDefault    bool
}

func (n *CreateOpClassStmt) Pos() int {
	return 0
}
