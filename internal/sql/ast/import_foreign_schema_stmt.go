package ast

type ImportForeignSchemaStmt struct {
	Tag NodeTag[ImportForeignSchemaStmt] `json:"tag"`

	ServerName   *string `json:",omitempty"`
	RemoteSchema *string `json:",omitempty"`
	LocalSchema  *string `json:",omitempty"`
	ListType     ImportForeignSchemaType
	TableList    *List `json:",omitempty"`
	Options      *List `json:",omitempty"`
}

func (n *ImportForeignSchemaStmt) Pos() int {
	return 0
}
