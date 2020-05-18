package ast

type AlterTableStmt struct {
	Table *TableName
	Cmds  *List
	// MissingOk bool
}

func (n *AlterTableStmt) Pos() int {
	return 0
}
