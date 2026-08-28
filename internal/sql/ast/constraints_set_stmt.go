package ast

type ConstraintsSetStmt struct {
	Tag NodeTag[ConstraintsSetStmt] `json:"tag"`

	Constraints *List `json:"constraints,omitempty"`
	Deferred    bool  `json:"deferred"`
}

func (n *ConstraintsSetStmt) Pos() int {
	return 0
}
