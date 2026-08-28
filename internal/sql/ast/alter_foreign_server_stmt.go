package ast

type AlterForeignServerStmt struct {
	Tag NodeTag[AlterForeignServerStmt] `json:"tag"`

	Servername *string `json:"servername,omitempty"`
	Version    *string `json:"version,omitempty"`
	Options    *List   `json:"options,omitempty"`
	HasVersion bool    `json:"has_version"`
}

func (n *AlterForeignServerStmt) Pos() int {
	return 0
}
