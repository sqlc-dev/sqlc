package ast

type DropTypeStmt struct {
	Tag NodeTag[DropTypeStmt] `json:"tag"`

	IfExists bool        `json:"if_exists"`
	Types    []*TypeName `json:"types,omitempty"`
}

func (n *DropTypeStmt) Pos() int {
	return 0
}
