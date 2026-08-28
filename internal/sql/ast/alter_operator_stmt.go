package ast

type AlterOperatorStmt struct {
	Tag NodeTag[AlterOperatorStmt] `json:"tag"`

	Opername *ObjectWithArgs `json:"opername,omitempty"`
	Options  *List           `json:"options,omitempty"`
}

func (n *AlterOperatorStmt) Pos() int {
	return 0
}
