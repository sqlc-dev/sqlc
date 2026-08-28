package ast

type CreateForeignTableStmt struct {
	Tag NodeTag[CreateForeignTableStmt] `json:"tag"`

	Base       *CreateStmt
	Servername *string
	Options    *List
}

func (n *CreateForeignTableStmt) Pos() int {
	return 0
}
