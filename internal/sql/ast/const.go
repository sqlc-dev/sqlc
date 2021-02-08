package ast

type Const struct {
	Xpr         Node
	Consttype   Oid
	Consttypmod int32
	Constcollid Oid
	Constlen    int
	Constvalue  Datum
	Constisnull bool
	Constbyval  bool
	Location    int
}

func (n *Const) Pos() int {
	return n.Location
}
