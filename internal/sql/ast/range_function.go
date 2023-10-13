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
	buf.astFormat(n.Functions)
	if n.Ordinality {
		buf.WriteString(" WITH ORDINALITY ")
	}
	buf.astFormat(n.Alias)
}
