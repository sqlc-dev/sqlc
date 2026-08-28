package ast

type DropTypeStmt struct {
	Tag NodeTag[DropTypeStmt] `json:"tag"`

	IfExists bool
	Types    []*TypeName `json:",omitempty"`
}

func (n *DropTypeStmt) Pos() int {
	return 0
}
