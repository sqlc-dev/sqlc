package ast

import ()

type AlterTableStmt struct {
	Relation  *RangeVar
	Cmds      *List
	Relkind   ObjectType
	MissingOk bool
}

func (n *AlterTableStmt) Pos() int {
	return 0
}
