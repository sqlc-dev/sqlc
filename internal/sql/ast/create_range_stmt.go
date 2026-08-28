package ast

type CreateRangeStmt struct {
	Tag NodeTag[CreateRangeStmt] `json:"tag"`

	TypeName *List `json:",omitempty"`
	Params   *List `json:",omitempty"`
}

func (n *CreateRangeStmt) Pos() int {
	return 0
}
