package ast

type AlterTableSpaceOptionsStmt struct {
	Tablespacename *string
	Options        *List
	IsReset        bool
}

func (n *AlterTableSpaceOptionsStmt) Pos() int {
	return 0
}
