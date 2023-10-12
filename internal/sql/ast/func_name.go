package ast

type FuncName struct {
	Catalog string
	Schema  string
	Name    string
}

func (n *FuncName) Pos() int {
	return 0
}

func (n *FuncName) Format(buf *TrackedBuffer) {
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
