package ast

type RangeVar struct {
	Catalogname    *string
	Schemaname     *string
	Relname        *string
	Inh            bool
	Relpersistence byte
	Alias          *Alias
	Location       int
}

func (n *RangeVar) Pos() int {
	return n.Location
}

func (n *RangeVar) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Schemaname != nil {
		buf.WriteString(*n.Schemaname)
		buf.WriteString(".")
	}
	if n.Relname != nil {
		// TODO: What names need to be quoted
		if *n.Relname == "user" {
			buf.WriteString(`"`)
			buf.WriteString(*n.Relname)
			buf.WriteString(`"`)
		} else {
			buf.WriteString(*n.Relname)
		}
	}
	if n.Alias != nil {
		buf.WriteString(" ")
		buf.astFormat(n.Alias)
	}
}
