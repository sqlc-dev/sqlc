package ast

type CompositeTypeStmt struct {
	TypeName *TypeName
	Coldefs  *List
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
