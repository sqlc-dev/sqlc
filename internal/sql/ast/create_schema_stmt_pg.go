package ast

type CreateSchemaStmt_PG struct {
	Schemaname  *string
	Authrole    *RoleSpec
	SchemaElts  *List
	IfNotExists bool
}

func (n *CreateSchemaStmt_PG) Pos() int {
	return 0
}
