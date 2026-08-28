package ast

type ImportForeignSchemaStmt struct {
	Tag NodeTag[ImportForeignSchemaStmt] `json:"tag"`

	ServerName   *string                 `json:"server_name,omitempty"`
	RemoteSchema *string                 `json:"remote_schema,omitempty"`
	LocalSchema  *string                 `json:"local_schema,omitempty"`
	ListType     ImportForeignSchemaType `json:"list_type"`
	TableList    *List                   `json:"table_list,omitempty"`
	Options      *List                   `json:"options,omitempty"`
}

func (n *ImportForeignSchemaStmt) Pos() int {
	return 0
}
