package ast

type RecursiveFuncCall struct {
	Func           Node
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

func (n *RecursiveFuncCall) Pos() int {
	return n.Location
}

func (n *RecursiveFuncCall) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.astFormat(n.Func)
	buf.WriteString("(")
	if n.AggStar {
		buf.WriteString("*")
	} else {
		buf.astFormat(n.Args)
	}
	buf.WriteString(")")
}
