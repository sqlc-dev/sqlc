package ast

type WindowFunc struct {
	Tag NodeTag[WindowFunc] `json:"tag"`

	Xpr         Node `json:",omitempty"`
	Winfnoid    Oid
	Wintype     Oid
	Wincollid   Oid
	Inputcollid Oid
	Args        *List `json:",omitempty"`
	Aggfilter   Node  `json:",omitempty"`
	Winref      Index
	Winstar     bool
	Winagg      bool
	Location    int
}

func (n *WindowFunc) Pos() int {
	return n.Location
}
