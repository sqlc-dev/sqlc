package ast

type RangeFunction struct {
	Lateral    bool
	Ordinality bool
	IsRowsfrom bool
	Functions  *List
	Alias      *Alias
	Coldeflist *List
}

func (n *RangeFunction) Pos() int {
	return 0
}

func (n *RangeFunction) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Lateral {
		buf.WriteString("LATERAL ")
	}
	buf.astFormat(n.Functions)
	if n.Ordinality {
		buf.WriteString(" WITH ORDINALITY")
	}
	if n.Alias != nil {
		buf.WriteString(" AS ")
		buf.astFormat(n.Alias)
	}
}
