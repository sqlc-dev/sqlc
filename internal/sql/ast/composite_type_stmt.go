package ast

type CompositeTypeStmt struct {
	TypeName   *TypeName
	ColDefList *List
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
