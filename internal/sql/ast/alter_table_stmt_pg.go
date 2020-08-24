package ast

type AlterTableStmt_PG struct {
	Relation  *RangeVar
	Cmds      *List
	Relkind   ObjectType
	MissingOk bool
}

func (n *AlterTableStmt_PG) Pos() int {
	return 0
}
