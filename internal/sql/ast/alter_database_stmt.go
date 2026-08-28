package ast

type AlterDatabaseStmt struct {
	Tag NodeTag[AlterDatabaseStmt] `json:"tag"`

	Dbname  *string `json:"dbname,omitempty"`
	Options *List   `json:"options,omitempty"`
}

func (n *AlterDatabaseStmt) Pos() int {
	return 0
}
