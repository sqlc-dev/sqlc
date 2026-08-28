package ast

type DeallocateStmt struct {
	Tag NodeTag[DeallocateStmt] `json:"tag"`

	Name *string
}

func (n *DeallocateStmt) Pos() int {
	return 0
}
