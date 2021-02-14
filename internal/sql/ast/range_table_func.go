package ast

type RangeTableFunc struct {
	Lateral    bool
	Docexpr    Node
	Rowexpr    Node
	Namespaces *List
	Columns    *List
	Alias      *Alias
	Location   int
}

func (n *RangeTableFunc) Pos() int {
	return n.Location
}
