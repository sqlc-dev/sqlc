package ast

type ConstraintsSetStmt struct {
	Constraints *List
	Deferred    bool
}

func (n *ConstraintsSetStmt) Pos() int {
	return 0
}
