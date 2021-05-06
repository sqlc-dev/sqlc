package ast

type AlterForeignServerStmt struct {
	Servername *string
	Version    *string
	Options    *List
	HasVersion bool
}

func (n *AlterForeignServerStmt) Pos() int {
	return 0
}
