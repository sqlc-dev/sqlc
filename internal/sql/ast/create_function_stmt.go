package ast

type CreateFunctionStmt struct {
	Replace    bool
	Params     *List
	ReturnType *TypeName
	Func       *FuncName
	// TODO: Understand these two fields
	Options    *List
	WithClause *List
}

func (n *CreateFunctionStmt) Pos() int {
	return 0
}
