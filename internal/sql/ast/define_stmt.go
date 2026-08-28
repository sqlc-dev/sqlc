package ast

type DefineStmt struct {
	Tag NodeTag[DefineStmt] `json:"tag"`

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
