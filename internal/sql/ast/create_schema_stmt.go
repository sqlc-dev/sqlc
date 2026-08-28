package ast

type CreateSchemaStmt struct {
	Tag NodeTag[CreateSchemaStmt] `json:"tag"`

	Name        *string   `json:",omitempty"`
	SchemaElts  *List     `json:",omitempty"`
	Authrole    *RoleSpec `json:",omitempty"`
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
