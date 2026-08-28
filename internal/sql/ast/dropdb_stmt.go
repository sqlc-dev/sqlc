package ast

type DropdbStmt struct {
	Tag NodeTag[DropdbStmt] `json:"tag"`

	Dbname    *string `json:"dbname,omitempty"`
	MissingOk bool    `json:"missing_ok"`
}

func (n *DropdbStmt) Pos() int {
	return 0
}
