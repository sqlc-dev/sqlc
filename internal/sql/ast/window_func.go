package ast

type WindowFunc struct {
	Xpr         Node
	Winfnoid    Oid
	Wintype     Oid
	Wincollid   Oid
	Inputcollid Oid
	Args        *List
	Aggfilter   Node
	Winref      Index
	Winstar     bool
	Winagg      bool
	Location    int
}

func (n *WindowFunc) Pos() int {
	return n.Location
}
