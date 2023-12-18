package ast

type MergeWhenClause struct {
	Condition  Node
	TargetList *List
	Values     *List
}

func (n *MergeWhenClause) Pos() int {
	return n.Pos()
}
