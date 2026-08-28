package ast

type AlterForeignServerStmt struct {
	Tag NodeTag[AlterForeignServerStmt] `json:"tag"`

	Servername *string `json:",omitempty"`
	Version    *string `json:",omitempty"`
	Options    *List   `json:",omitempty"`
	HasVersion bool
}

func (n *AlterForeignServerStmt) Pos() int {
	return 0
}
