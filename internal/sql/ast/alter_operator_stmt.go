package ast

type AlterOperatorStmt struct {
	Tag NodeTag[AlterOperatorStmt] `json:"tag"`

	Opername *ObjectWithArgs `json:",omitempty"`
	Options  *List           `json:",omitempty"`
}

func (n *AlterOperatorStmt) Pos() int {
	return 0
}
