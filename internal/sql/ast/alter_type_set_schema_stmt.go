package ast

type AlterTypeSetSchemaStmt struct {
	Tag NodeTag[AlterTypeSetSchemaStmt] `json:"tag"`

	Type      *TypeName
	NewSchema *string
}

func (n *AlterTypeSetSchemaStmt) Pos() int {
	return 0
}
