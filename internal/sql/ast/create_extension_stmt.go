package ast

type CreateExtensionStmt struct {
	Extname     *string
	IfNotExists bool
	Options     *List
}

func (n *CreateExtensionStmt) Pos() int {
	return 0
}

func (n *CreateExtensionStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("CREATE EXTENSION ")
	if n.IfNotExists {
		buf.WriteString("IF NOT EXISTS ")
	}
	if n.Extname != nil {
		buf.WriteString(*n.Extname)
	}
}
