package ast

type AlterTSDictionaryStmt struct {
	Tag NodeTag[AlterTSDictionaryStmt] `json:"tag"`

	Dictname *List
	Options  *List
}

func (n *AlterTSDictionaryStmt) Pos() int {
	return 0
}
