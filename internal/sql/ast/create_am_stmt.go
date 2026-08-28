package ast

type CreateAmStmt struct {
	Tag NodeTag[CreateAmStmt] `json:"tag"`

	Amname      *string `json:",omitempty"`
	HandlerName *List   `json:",omitempty"`
	Amtype      byte
}

func (n *CreateAmStmt) Pos() int {
	return 0
}
