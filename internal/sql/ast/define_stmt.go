package ast

type DefineStmt struct {
	Tag NodeTag[DefineStmt] `json:"tag"`

	Kind        ObjectType `json:"kind"`
	Oldstyle    bool       `json:"oldstyle"`
	Defnames    *List      `json:"defnames,omitempty"`
	Args        *List      `json:"args,omitempty"`
	Definition  *List      `json:"definition,omitempty"`
	IfNotExists bool       `json:"if_not_exists"`
}

func (n *DefineStmt) Pos() int {
	return 0
}
