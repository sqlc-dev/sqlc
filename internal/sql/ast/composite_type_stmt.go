package ast

type CompositeTypeStmt struct {
	TypeName *TypeName
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
