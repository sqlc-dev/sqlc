package ast

type AlterFunctionStmt struct {
	Tag NodeTag[AlterFunctionStmt] `json:"tag"`

	Func    *ObjectWithArgs `json:"func,omitempty"`
	Actions *List           `json:"actions,omitempty"`
}

func (n *AlterFunctionStmt) Pos() int {
	return 0
}
