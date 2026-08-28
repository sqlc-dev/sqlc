package ast

type ArrayExpr struct {
	Tag NodeTag[ArrayExpr] `json:"tag"`

	Xpr           Node `json:",omitempty"`
	ArrayTypeid   Oid
	ArrayCollid   Oid
	ElementTypeid Oid
	Elements      *List `json:",omitempty"`
	Multidims     bool
	Location      int
}

func (n *ArrayExpr) Pos() int {
	return n.Location
}
