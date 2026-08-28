package ast

type DefineStmt struct {
	Tag NodeTag[DefineStmt] `json:"tag"`

	Kind        ObjectType
	Oldstyle    bool
	Defnames    *List `json:",omitempty"`
	Args        *List `json:",omitempty"`
	Definition  *List `json:",omitempty"`
	IfNotExists bool
}

func (n *DefineStmt) Pos() int {
	return 0
}
