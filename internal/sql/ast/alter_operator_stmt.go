package ast

type AlterOperatorStmt struct {
	Opername *ObjectWithArgs
	Options  *List
}

func (n *AlterOperatorStmt) Pos() int {
	return 0
}
