package ast

type AlterDatabaseStmt struct {
	Tag NodeTag[AlterDatabaseStmt] `json:"tag"`

	Dbname  *string
	Options *List
}

func (n *AlterDatabaseStmt) Pos() int {
	return 0
}
