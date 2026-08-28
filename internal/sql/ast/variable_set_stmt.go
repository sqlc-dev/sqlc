package ast

type VariableSetStmt struct {
	Tag NodeTag[VariableSetStmt] `json:"tag"`

	Kind    VariableSetKind `json:"kind"`
	Name    *string         `json:"name,omitempty"`
	Args    *List           `json:"args,omitempty"`
	IsLocal bool            `json:"is_local"`
}

func (n *VariableSetStmt) Pos() int {
	return 0
}
