package ast

type RangeTableFunc struct {
	Tag NodeTag[RangeTableFunc] `json:"tag"`

	Lateral    bool   `json:"lateral"`
	Docexpr    Node   `json:"docexpr,omitempty"`
	Rowexpr    Node   `json:"rowexpr,omitempty"`
	Namespaces *List  `json:"namespaces,omitempty"`
	Columns    *List  `json:"columns,omitempty"`
	Alias      *Alias `json:"alias,omitempty"`
	Location   int    `json:"location"`
}

func (n *RangeTableFunc) Pos() int {
	return n.Location
}
