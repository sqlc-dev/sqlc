package ast

type AlterTypeSetSchemaStmt struct {
	Type      *TypeName
	NewSchema *string
}

func (n *AlterTypeSetSchemaStmt) Pos() int {
	return 0
}
