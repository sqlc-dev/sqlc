package ast

type AlterPublicationStmt struct {
	Tag NodeTag[AlterPublicationStmt] `json:"tag"`

	Pubname      *string `json:",omitempty"`
	Options      *List   `json:",omitempty"`
	Tables       *List   `json:",omitempty"`
	ForAllTables bool
	TableAction  DefElemAction
}

func (n *AlterPublicationStmt) Pos() int {
	return 0
}
