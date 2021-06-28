package ast

type AlterSeqStmt struct {
	Sequence    *RangeVar
	Options     *List
	ForIdentity bool
	MissingOk   bool
}

func (n *AlterSeqStmt) Pos() int {
	return 0
}
