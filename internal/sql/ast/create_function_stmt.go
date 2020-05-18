package ast

type CreateFunctionStmt struct {
	Replace    bool
	Params     []*FuncParam
	ReturnType *TypeName
	Func       *FuncName
}

func (n *CreateFunctionStmt) Pos() int {
	return 0
}
