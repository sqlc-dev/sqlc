package ast

type Var struct {
	Tag NodeTag[Var] `json:"tag"`

	Xpr         Node       `json:"xpr,omitempty"`
	Varno       Index      `json:"varno"`
	Varattno    AttrNumber `json:"varattno"`
	Vartype     Oid        `json:"vartype"`
	Vartypmod   int32      `json:"vartypmod"`
	Varcollid   Oid        `json:"varcollid"`
	Varlevelsup Index      `json:"varlevelsup"`
	Varnoold    Index      `json:"varnoold"`
	Varoattno   AttrNumber `json:"varoattno"`
	Location    int        `json:"location"`
}

func (n *Var) Pos() int {
	return n.Location
}
