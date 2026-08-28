package ast

type WindowFunc struct {
	Tag NodeTag[WindowFunc] `json:"tag"`

	Xpr         Node  `json:"xpr,omitempty"`
	Winfnoid    Oid   `json:"winfnoid"`
	Wintype     Oid   `json:"wintype"`
	Wincollid   Oid   `json:"wincollid"`
	Inputcollid Oid   `json:"inputcollid"`
	Args        *List `json:"args,omitempty"`
	Aggfilter   Node  `json:"aggfilter,omitempty"`
	Winref      Index `json:"winref"`
	Winstar     bool  `json:"winstar"`
	Winagg      bool  `json:"winagg"`
	Location    int   `json:"location"`
}

func (n *WindowFunc) Pos() int {
	return n.Location
}
