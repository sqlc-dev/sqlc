package ast

type CompositeTypeStmt struct {
	Tag NodeTag[CompositeTypeStmt] `json:"tag"`

	TypeName *TypeName
}

func (n *CompositeTypeStmt) Pos() int {
	return 0
}
