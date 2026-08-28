package ast

type AlterTableSpaceOptionsStmt struct {
	Tag NodeTag[AlterTableSpaceOptionsStmt] `json:"tag"`

	Tablespacename *string `json:",omitempty"`
	Options        *List   `json:",omitempty"`
	IsReset        bool
}

func (n *AlterTableSpaceOptionsStmt) Pos() int {
	return 0
}
