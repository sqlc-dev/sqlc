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
	if n.AggStar {
		buf.WriteString("*")
	} else {
		buf.astFormat(n.Args)
	}
	buf.WriteString(")")
}
