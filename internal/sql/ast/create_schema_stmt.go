package ast

type CreateSchemaStmt struct {
	Tag NodeTag[CreateSchemaStmt] `json:"tag"`

	Name        *string
	SchemaElts  *List
	Authrole    *RoleSpec
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
