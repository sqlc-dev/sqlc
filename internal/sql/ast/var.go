package ast

type Var struct {
	Xpr         Node
	Varno       Index
	Varattno    AttrNumber
	Vartype     Oid
	Vartypmod   int32
	Varcollid   Oid
	Varlevelsup Index
	Varnoold    Index
	Varoattno   AttrNumber
	Location    int
}

func (n *Var) Pos() int {
	return n.Location
}
