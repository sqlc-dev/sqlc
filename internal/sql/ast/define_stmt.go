package ast

type DefineStmt struct {
	Kind        ObjectType
	Oldstyle    bool
	Defnames    *List
	Args        *List
	Definition  *List
	IfNotExists bool
}

func (n *DefineStmt) Pos() int {
	return 0
}
