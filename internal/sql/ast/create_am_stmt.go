package ast

type CreateAmStmt struct {
	Amname      *string
	HandlerName *List
	Amtype      byte
}

func (n *CreateAmStmt) Pos() int {
	return 0
}
