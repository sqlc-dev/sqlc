package ast

type AlterTSDictionaryStmt struct {
	Tag NodeTag[AlterTSDictionaryStmt] `json:"tag"`

	Dictname *List `json:"dictname,omitempty"`
	Options  *List `json:"options,omitempty"`
}

func (n *AlterTSDictionaryStmt) Pos() int {
	return 0
}
