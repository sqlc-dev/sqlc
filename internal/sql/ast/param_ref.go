package ast

import "fmt"

type ParamRef struct {
	Number   int
	Location int
	Dollar   bool
}

func (n *ParamRef) Pos() int {
	return n.Location
}

func (n *ParamRef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	fmt.Fprintf(buf, "$%d", n.Number)
}
