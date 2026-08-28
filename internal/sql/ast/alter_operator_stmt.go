package ast

type AlterOperatorStmt struct {
	Tag NodeTag[AlterOperatorStmt] `json:"tag"`

	Opername *ObjectWithArgs
	Options  *List
}

func (n *AlterOperatorStmt) Pos() int {
	return 0
}
