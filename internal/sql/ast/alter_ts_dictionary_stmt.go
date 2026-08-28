package ast

type AlterTSDictionaryStmt struct {
	Tag NodeTag[AlterTSDictionaryStmt] `json:"tag"`

	Dictname *List `json:",omitempty"`
	Options  *List `json:",omitempty"`
}

func (n *AlterTSDictionaryStmt) Pos() int {
	return 0
}
