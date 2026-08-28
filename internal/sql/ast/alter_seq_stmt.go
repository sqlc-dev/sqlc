package ast

type AlterSeqStmt struct {
	Tag NodeTag[AlterSeqStmt] `json:"tag"`

	Sequence    *RangeVar `json:"sequence,omitempty"`
	Options     *List     `json:"options,omitempty"`
	ForIdentity bool      `json:"for_identity"`
	MissingOk   bool      `json:"missing_ok"`
}

func (n *AlterSeqStmt) Pos() int {
	return 0
}
