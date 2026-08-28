package ast

type AlterSeqStmt struct {
	Tag NodeTag[AlterSeqStmt] `json:"tag"`

	Sequence    *RangeVar
	Options     *List
	ForIdentity bool
	MissingOk   bool
}

func (n *AlterSeqStmt) Pos() int {
	return 0
}
