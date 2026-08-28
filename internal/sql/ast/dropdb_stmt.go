package ast

type DropdbStmt struct {
	Tag NodeTag[DropdbStmt] `json:"tag"`

	Dbname    *string `json:",omitempty"`
	MissingOk bool
}

func (n *DropdbStmt) Pos() int {
	return 0
}
