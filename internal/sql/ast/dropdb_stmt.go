package ast

type DropdbStmt struct {
	Tag NodeTag[DropdbStmt] `json:"tag"`

	Dbname    *string
	MissingOk bool
}

func (n *DropdbStmt) Pos() int {
	return 0
}
