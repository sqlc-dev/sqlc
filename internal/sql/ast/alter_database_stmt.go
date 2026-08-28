package ast

type AlterDatabaseStmt struct {
	Tag NodeTag[AlterDatabaseStmt] `json:"tag"`

	Dbname  *string `json:",omitempty"`
	Options *List   `json:",omitempty"`
}

func (n *AlterDatabaseStmt) Pos() int {
	return 0
}
