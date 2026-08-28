package ast

type Param struct {
	Tag NodeTag[Param] `json:"tag"`

	Xpr         Node `json:",omitempty"`
	Paramkind   ParamKind
	Paramid     int
	Paramtype   Oid
	Paramtypmod int32
	Paramcollid Oid
	Location    int
}

func (n *Param) Pos() int {
	return n.Location
}
