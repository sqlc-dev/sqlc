package ast

type DefElem struct {
	Defnamespace *string
	Defname      *string
	Arg          Node
	Defaction    DefElemAction
	Location     int
}

func (n *DefElem) Pos() int {
	return n.Location
}
