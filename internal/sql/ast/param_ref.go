package ast

type ParamRef struct {
	Number   int
	Location int
}

func (n *ParamRef) Pos() int {
	return n.Location
}
