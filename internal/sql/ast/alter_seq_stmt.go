package ast

type AlterSeqStmt struct {
	Tag NodeTag[AlterSeqStmt] `json:"tag"`

	Sequence    *RangeVar `json:",omitempty"`
	Options     *List     `json:",omitempty"`
	ForIdentity bool
	MissingOk   bool
}

func (n *AlterSeqStmt) Pos() int {
	return 0
}
