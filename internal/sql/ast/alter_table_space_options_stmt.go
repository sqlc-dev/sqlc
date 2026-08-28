package ast

type AlterTableSpaceOptionsStmt struct {
	Tag NodeTag[AlterTableSpaceOptionsStmt] `json:"tag"`

	Tablespacename *string
	Options        *List
	IsReset        bool
}

func (n *AlterTableSpaceOptionsStmt) Pos() int {
	return 0
}
