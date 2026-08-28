package ast

type CreatePublicationStmt struct {
	Tag NodeTag[CreatePublicationStmt] `json:"tag"`

	Pubname      *string
	Options      *List
	Tables       *List
	ForAllTables bool
}

func (n *CreatePublicationStmt) Pos() int {
	return 0
}
