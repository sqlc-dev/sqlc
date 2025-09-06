package ast

import "fmt"

type ParamRef struct {
	Number   int
	Location int
	Dollar   bool

	// YDB specific 
	Plike bool
}

func (n *ParamRef) Pos() int {
	return n.Location
}

func (n *ParamRef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	if n.Plike {
		fmt.Fprintf(buf, "$p%d", n.Number)
		return
	}
	fmt.Fprintf(buf, "$%d", n.Number)
}
