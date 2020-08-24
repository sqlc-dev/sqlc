package ast

import ()

type CreateSchemaStmt struct {
	Schemaname  *string
	Authrole    *RoleSpec
	SchemaElts  *List
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
