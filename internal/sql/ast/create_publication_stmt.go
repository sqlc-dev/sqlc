package ast

type CreatePublicationStmt struct {
	Tag NodeTag[CreatePublicationStmt] `json:"tag"`

	Pubname      *string `json:"pubname,omitempty"`
	Options      *List   `json:"options,omitempty"`
	Tables       *List   `json:"tables,omitempty"`
	ForAllTables bool    `json:"for_all_tables"`
}

func (n *CreatePublicationStmt) Pos() int {
	return 0
}
