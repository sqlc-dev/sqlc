package ast

type Const struct {
	Tag NodeTag[Const] `json:"tag"`

	Xpr         Node  `json:"xpr,omitempty"`
	Consttype   Oid   `json:"consttype"`
	Consttypmod int32 `json:"consttypmod"`
	Constcollid Oid   `json:"constcollid"`
	Constlen    int   `json:"constlen"`
	Constvalue  Datum `json:"constvalue"`
	Constisnull bool  `json:"constisnull"`
	Constbyval  bool  `json:"constbyval"`
	Location    int   `json:"location"`
}

func (n *Const) Pos() int {
	return n.Location
}
