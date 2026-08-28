package ast

type AlterForeignServerStmt struct {
	Tag NodeTag[AlterForeignServerStmt] `json:"tag"`

	Servername *string
	Version    *string
	Options    *List
	HasVersion bool
}

func (n *AlterForeignServerStmt) Pos() int {
	return 0
}
