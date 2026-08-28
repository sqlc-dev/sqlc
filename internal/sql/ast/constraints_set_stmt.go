package ast

type ConstraintsSetStmt struct {
	Tag NodeTag[ConstraintsSetStmt] `json:"tag"`

	Constraints *List `json:",omitempty"`
	Deferred    bool
}

func (n *ConstraintsSetStmt) Pos() int {
	return 0
}
