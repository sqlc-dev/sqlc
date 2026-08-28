package ast

type ConstraintsSetStmt struct {
	Tag NodeTag[ConstraintsSetStmt] `json:"tag"`

	Constraints *List
	Deferred    bool
}

func (n *ConstraintsSetStmt) Pos() int {
	return 0
}
