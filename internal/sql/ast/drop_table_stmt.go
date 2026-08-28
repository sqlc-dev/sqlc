package ast

type DropTableStmt struct {
	Tag NodeTag[DropTableStmt] `json:"tag"`

	IfExists bool
	Tables   []*TableName `json:",omitempty"`
}

func (n *DropTableStmt) Pos() int {
	return 0
}
