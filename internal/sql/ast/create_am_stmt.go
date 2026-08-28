package ast

type CreateAmStmt struct {
	Tag NodeTag[CreateAmStmt] `json:"tag"`

	Amname      *string `json:"amname,omitempty"`
	HandlerName *List   `json:"handler_name,omitempty"`
	Amtype      byte    `json:"amtype"`
}

func (n *CreateAmStmt) Pos() int {
	return 0
}
