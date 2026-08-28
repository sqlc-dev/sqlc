package ast

type DropTableStmt struct {
	Tag NodeTag[DropTableStmt] `json:"tag"`

	IfExists bool
	Tables   []*TableName
}

func (n *DropTableStmt) Pos() int {
	return 0
}
