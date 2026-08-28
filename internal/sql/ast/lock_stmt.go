package ast

type LockStmt struct {
	Tag NodeTag[LockStmt] `json:"tag"`

	Relations *List `json:",omitempty"`
	Mode      int
	Nowait    bool
}

func (n *LockStmt) Pos() int {
	return 0
}
