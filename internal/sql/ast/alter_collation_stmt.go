package ast

type AlterCollationStmt struct {
	Collname *List
}

func (n *AlterCollationStmt) Pos() int {
	return 0
}
