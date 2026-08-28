package ast

type CreateAmStmt struct {
	Tag NodeTag[CreateAmStmt] `json:"tag"`

	Amname      *string
	HandlerName *List
	Amtype      byte
}

func (n *CreateAmStmt) Pos() int {
	return 0
}
