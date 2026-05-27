package ast

type DropTableStmt struct {
	IfExists bool
	Tables   []*TableName
	Behavior DropBehavior
}

func (n *DropTableStmt) Pos() int {
	return 0
}
