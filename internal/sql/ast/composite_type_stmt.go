package ast

type CompositeTypeStmt struct {
	Tag NodeTag[CompositeTypeStmt] `json:"tag"`

	TypeName *TypeName `json:",omitempty"`
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
