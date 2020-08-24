package ast

type CreateSchemaStmt struct {
	Name        *string
	SchemaElts  *List
	Authrole    *RoleSpec
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
