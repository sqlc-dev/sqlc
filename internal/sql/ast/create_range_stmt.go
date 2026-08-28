package ast

type CreateRangeStmt struct {
	Tag NodeTag[CreateRangeStmt] `json:"tag"`

	TypeName *List
	Params   *List
}

func (n *CreateRangeStmt) Pos() int {
	return 0
}
