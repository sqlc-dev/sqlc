package ast

type ScalarArrayOpExpr struct {
	Xpr         Node
	Opno        Oid
	UseOr       bool
	Inputcollid Oid
	Args        *List
	Location    int
}

func (n *ScalarArrayOpExpr) Pos() int {
	return n.Location
}
