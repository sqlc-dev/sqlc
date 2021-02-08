package ast

type CreatedbStmt struct {
	Dbname  *string
	Options *List
}

func (n *CreatedbStmt) Pos() int {
	return 0
}
