package ast

type LockStmt struct {
	Tag NodeTag[LockStmt] `json:"tag"`

	Relations *List `json:"relations,omitempty"`
	Mode      int   `json:"mode"`
	Nowait    bool  `json:"nowait"`
}

func (n *LockStmt) Pos() int {
	return 0
}
