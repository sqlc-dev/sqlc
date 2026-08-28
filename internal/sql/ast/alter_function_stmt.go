package ast

type AlterFunctionStmt struct {
	Tag NodeTag[AlterFunctionStmt] `json:"tag"`

	Func    *ObjectWithArgs
	Actions *List
}

func (n *AlterFunctionStmt) Pos() int {
	return 0
}
