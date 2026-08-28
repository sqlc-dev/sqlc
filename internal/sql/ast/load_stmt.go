package ast

type LoadStmt struct {
	Tag NodeTag[LoadStmt] `json:"tag"`

	Filename *string
}

func (n *LoadStmt) Pos() int {
	return 0
}
