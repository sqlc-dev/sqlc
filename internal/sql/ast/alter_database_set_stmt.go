package ast

type AlterDatabaseSetStmt struct {
	Tag NodeTag[AlterDatabaseSetStmt] `json:"tag"`

	Dbname  *string          `json:"dbname,omitempty"`
	Setstmt *VariableSetStmt `json:"setstmt,omitempty"`
}

func (n *AlterDatabaseSetStmt) Pos() int {
	return 0
}
