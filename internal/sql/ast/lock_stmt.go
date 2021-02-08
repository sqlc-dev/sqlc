package ast

type LockStmt struct {
	Relations *List
	Mode      int
	Nowait    bool
}

func (n *LockStmt) Pos() int {
	return 0
}
