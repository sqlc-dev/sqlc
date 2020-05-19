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
