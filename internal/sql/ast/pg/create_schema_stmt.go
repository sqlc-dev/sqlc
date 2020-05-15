package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateSchemaStmt struct {
	Schemaname  *string
	Authrole    *RoleSpec
	SchemaElts  *ast.List
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
