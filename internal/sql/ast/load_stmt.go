package ast

type LoadStmt struct {
	Tag NodeTag[LoadStmt] `json:"tag"`

	Filename *string `json:"filename,omitempty"`
}

func (n *LoadStmt) Pos() int {
	return 0
}
