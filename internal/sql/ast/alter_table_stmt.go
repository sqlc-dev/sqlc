package ast

type AlterTableStmt struct {
	// TODO: Only TableName or Relation should be defined
	Relation  *RangeVar
	Table     *TableName
	Cmds      *List
	MissingOk bool
	Relkind   ObjectType
}

func (n *AlterTableStmt) Pos() int {
	return 0
}
