package ast

type WindowFunc struct {
	Tag NodeTag[WindowFunc] `json:"tag"`

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
