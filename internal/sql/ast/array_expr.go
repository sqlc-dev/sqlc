package ast

type ArrayExpr struct {
	Tag NodeTag[ArrayExpr] `json:"tag"`

	Xpr           Node
	ArrayTypeid   Oid
	ArrayCollid   Oid
	ElementTypeid Oid
	Elements      *List
	Multidims     bool
	Location      int
}

func (n *ArrayExpr) Pos() int {
	return n.Location
}
