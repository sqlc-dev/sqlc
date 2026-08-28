package ast

type SetOperationStmt struct {
	Tag NodeTag[SetOperationStmt] `json:"tag"`

	Op            SetOperation
	All           bool
	Larg          Node  `json:",omitempty"`
	Rarg          Node  `json:",omitempty"`
	ColTypes      *List `json:",omitempty"`
	ColTypmods    *List `json:",omitempty"`
	ColCollations *List `json:",omitempty"`
	GroupClauses  *List `json:",omitempty"`
}

func (n *SetOperationStmt) Pos() int {
	return 0
}
