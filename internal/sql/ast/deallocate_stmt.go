package ast

type DeallocateStmt struct {
	Tag NodeTag[DeallocateStmt] `json:"tag"`

	Name *string `json:"name,omitempty"`
}

func (n *DeallocateStmt) Pos() int {
	return 0
}
