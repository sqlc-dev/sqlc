package ast

type ImportForeignSchemaStmt struct {
	ServerName   *string
	RemoteSchema *string
	LocalSchema  *string
	ListType     ImportForeignSchemaType
	TableList    *List
	Options      *List
}

func (n *ImportForeignSchemaStmt) Pos() int {
	return 0
}
