package ast

type Param struct {
	Tag NodeTag[Param] `json:"tag"`

	Xpr         Node      `json:"xpr,omitempty"`
	Paramkind   ParamKind `json:"paramkind"`
	Paramid     int       `json:"paramid"`
	Paramtype   Oid       `json:"paramtype"`
	Paramtypmod int32     `json:"paramtypmod"`
	Paramcollid Oid       `json:"paramcollid"`
	Location    int       `json:"location"`
}

func (n *Param) Pos() int {
	return n.Location
}
