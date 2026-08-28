package ast

type CreatePublicationStmt struct {
	Tag NodeTag[CreatePublicationStmt] `json:"tag"`

	Pubname      *string `json:",omitempty"`
	Options      *List   `json:",omitempty"`
	Tables       *List   `json:",omitempty"`
	ForAllTables bool
}

func (n *CreatePublicationStmt) Pos() int {
	return 0
}
