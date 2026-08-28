package ast

type CreateRangeStmt struct {
	Tag NodeTag[CreateRangeStmt] `json:"tag"`

	TypeName *List `json:"type_name,omitempty"`
	Params   *List `json:"params,omitempty"`
}

func (n *CreateRangeStmt) Pos() int {
	return 0
}
