package ast

type DropSchemaStmt struct {
	Schemas   []*String
	MissingOk bool
}

func (n *DropSchemaStmt) Pos() int {
	return 0
}
