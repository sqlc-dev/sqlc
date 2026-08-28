package ast

type CreateSchemaStmt struct {
	Tag NodeTag[CreateSchemaStmt] `json:"tag"`

	Name        *string   `json:"name,omitempty"`
	SchemaElts  *List     `json:"schema_elts,omitempty"`
	Authrole    *RoleSpec `json:"authrole,omitempty"`
	IfNotExists bool      `json:"if_not_exists"`
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
