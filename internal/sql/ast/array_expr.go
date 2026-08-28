package ast

type ArrayExpr struct {
	Tag NodeTag[ArrayExpr] `json:"tag"`

	Xpr           Node  `json:"xpr,omitempty"`
	ArrayTypeid   Oid   `json:"array_typeid"`
	ArrayCollid   Oid   `json:"array_collid"`
	ElementTypeid Oid   `json:"element_typeid"`
	Elements      *List `json:"elements,omitempty"`
	Multidims     bool  `json:"multidims"`
	Location      int   `json:"location"`
}

func (n *ArrayExpr) Pos() int {
	return n.Location
}
