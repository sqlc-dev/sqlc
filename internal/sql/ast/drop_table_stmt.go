package ast

type DropTableStmt struct {
	IfExists bool
	Tables   []*TableName
}

func (n *DropTableStmt) Pos() int {
	return 0
}
