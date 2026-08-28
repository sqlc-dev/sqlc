package ast

type AlterDatabaseSetStmt struct {
	Tag NodeTag[AlterDatabaseSetStmt] `json:"tag"`

	Dbname  *string          `json:",omitempty"`
	Setstmt *VariableSetStmt `json:",omitempty"`
}

func (n *AlterDatabaseSetStmt) Pos() int {
	return 0
}
