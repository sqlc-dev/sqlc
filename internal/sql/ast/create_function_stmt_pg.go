package ast

type CreateFunctionStmt_PG struct {
	Replace    bool
	Funcname   *List
	Parameters *List
	ReturnType *TypeName
	Options    *List
	WithClause *List
}

func (n *CreateFunctionStmt_PG) Pos() int {
	return 0
}
