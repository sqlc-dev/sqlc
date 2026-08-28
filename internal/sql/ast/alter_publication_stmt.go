package ast

type AlterPublicationStmt struct {
	Tag NodeTag[AlterPublicationStmt] `json:"tag"`

	Pubname      *string       `json:"pubname,omitempty"`
	Options      *List         `json:"options,omitempty"`
	Tables       *List         `json:"tables,omitempty"`
	ForAllTables bool          `json:"for_all_tables"`
	TableAction  DefElemAction `json:"table_action"`
}

func (n *AlterPublicationStmt) Pos() int {
	return 0
}
