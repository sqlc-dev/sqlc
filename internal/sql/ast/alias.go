package ast

type Alias struct {
	Aliasname *string
	Colnames  *List
}

func (n *Alias) Pos() int {
	return 0
}

func (n *Alias) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Aliasname != nil {
		buf.WriteString(*n.Aliasname)
	}
	if items(n.Colnames) {
		buf.WriteString("(")
		buf.astFormat((n.Colnames))
		buf.WriteString(")")
	}
}
