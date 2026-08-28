package ast

type AlterTypeSetSchemaStmt struct {
	Tag NodeTag[AlterTypeSetSchemaStmt] `json:"tag"`

	Type      *TypeName `json:",omitempty"`
	NewSchema *string   `json:",omitempty"`
}

func (n *AlterTypeSetSchemaStmt) Pos() int {
	return 0
}
