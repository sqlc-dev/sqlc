package ast

type TableName struct {
	Catalog      string
	Schema       string
	Name         string // table name, maybe alias name, maybe original name
	OriginalName string // table original name
}

func (n *TableName) Pos() int {
	return 0
}

func (n *TableName) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Schema != "" {
		buf.WriteString(n.Schema)
		buf.WriteString(".")
	}
	if n.Name != "" {
		buf.WriteString(n.Name)
	}
}
