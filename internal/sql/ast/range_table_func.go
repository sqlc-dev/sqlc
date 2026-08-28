package ast

type RangeTableFunc struct {
	Tag NodeTag[RangeTableFunc] `json:"tag"`

	Lateral    bool
	Docexpr    Node   `json:",omitempty"`
	Rowexpr    Node   `json:",omitempty"`
	Namespaces *List  `json:",omitempty"`
	Columns    *List  `json:",omitempty"`
	Alias      *Alias `json:",omitempty"`
	Location   int
}

func (n *RangeTableFunc) Pos() int {
	return n.Location
}
