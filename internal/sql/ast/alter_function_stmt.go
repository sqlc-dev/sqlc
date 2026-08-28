package ast

type AlterFunctionStmt struct {
	Tag NodeTag[AlterFunctionStmt] `json:"tag"`

	Func    *ObjectWithArgs `json:",omitempty"`
	Actions *List           `json:",omitempty"`
}

func (n *AlterFunctionStmt) Pos() int {
	return 0
}
