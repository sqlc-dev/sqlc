package ast

type AlterTableSpaceOptionsStmt struct {
	Tag NodeTag[AlterTableSpaceOptionsStmt] `json:"tag"`

	Tablespacename *string `json:"tablespacename,omitempty"`
	Options        *List   `json:"options,omitempty"`
	IsReset        bool    `json:"is_reset"`
}

func (n *AlterTableSpaceOptionsStmt) Pos() int {
	return 0
}
