package ast

type VariableSetStmt struct {
	Tag NodeTag[VariableSetStmt] `json:"tag"`

	Kind    VariableSetKind
	Name    *string `json:",omitempty"`
	Args    *List   `json:",omitempty"`
	IsLocal bool
}

func (n *VariableSetStmt) Pos() int {
	return 0
}
