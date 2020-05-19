package ast

type CreateSchemaStmt struct {
	Name        *string
	IfNotExists bool
}

func (n *CreateSchemaStmt) Pos() int {
	return 0
}
