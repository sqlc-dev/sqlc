package ast

type AlterDatabaseSetStmt struct {
	Tag NodeTag[AlterDatabaseSetStmt] `json:"tag"`

	Dbname  *string
	Setstmt *VariableSetStmt
}

func (n *AlterDatabaseSetStmt) Pos() int {
	return 0
}
