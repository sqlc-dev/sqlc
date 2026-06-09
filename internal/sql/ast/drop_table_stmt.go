package ast

type DropTableStmt struct {
	Behavior  DropBehavior
	IfExists  bool
	Tables    []*TableName
}

func (n *DropTableStmt) Pos() int {
	return 0
}
