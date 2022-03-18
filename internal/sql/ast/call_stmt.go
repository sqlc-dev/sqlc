package ast

type CallStmt struct {
	FuncCall FuncCall
}

func (n *CallStmt) Pos() int {
	return n.FuncCall.Pos()
}
