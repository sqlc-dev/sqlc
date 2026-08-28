package ast

type CreatedbStmt struct {
	Tag NodeTag[CreatedbStmt] `json:"tag"`

	Dbname  *string
	Options *List
}

func (n *CreatedbStmt) Pos() int {
	return 0
}
