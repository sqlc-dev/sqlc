package ast

type FuncCall struct {
	Func           *FuncName
	Funcname       *List
	Args           *List
	AggOrder       *List
	AggFilter      Node
	AggWithinGroup bool
	AggStar        bool
	AggDistinct    bool
	FuncVariadic   bool
	Over           *WindowDef
	Location       int
}

func (n *FuncCall) Pos() int {
	return n.Location
}

func (n *FuncCall) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Func)
	buf.WriteString("(")
	if n.AggDistinct {
		buf.WriteString("DISTINCT ")
	}
	if n.AggStar {
		buf.WriteString("*")
	} else {
		buf.astFormat(n.Args)
	}
	if items(n.AggOrder) {
		buf.WriteString(" ORDER BY ")
		buf.join(n.AggOrder, ", ")
	}
	buf.WriteString(")")
	if set(n.AggFilter) {
		buf.WriteString(" FILTER (WHERE ")
		buf.astFormat(n.AggFilter)
		buf.WriteString(")")
	}
	if n.Over != nil {
		buf.WriteString(" OVER ")
		buf.astFormat(n.Over)
	}
}
