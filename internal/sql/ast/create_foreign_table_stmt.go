package ast

type CreateForeignTableStmt struct {
	Base       *CreateStmt
	Servername *string
	Options    *List
}

func (n *CreateForeignTableStmt) Pos() int {
	return 0
}
