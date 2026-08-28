package ast

type AlterTypeSetSchemaStmt struct {
	Tag NodeTag[AlterTypeSetSchemaStmt] `json:"tag"`

	Type      *TypeName `json:"type,omitempty"`
	NewSchema *string   `json:"new_schema,omitempty"`
}

func (n *AlterTypeSetSchemaStmt) Pos() int {
	return 0
}
