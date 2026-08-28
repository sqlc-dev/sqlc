package ast

type AlterSystemStmt struct {
	Tag NodeTag[AlterSystemStmt] `json:"tag"`

	Setstmt *VariableSetStmt `json:",omitempty"`
}

func (n *AlterSystemStmt) Pos() int {
	return 0
}
