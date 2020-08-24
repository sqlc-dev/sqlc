package ast

type AlterDatabaseSetStmt struct {
	Dbname  *string
	Setstmt *VariableSetStmt
}

func (n *AlterDatabaseSetStmt) Pos() int {
	return 0
}
