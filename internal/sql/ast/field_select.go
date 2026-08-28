package ast

type FieldSelect struct {
	Tag NodeTag[FieldSelect] `json:"tag"`

	Xpr          Node       `json:"xpr,omitempty"`
	Arg          Node       `json:"arg,omitempty"`
	Fieldnum     AttrNumber `json:"fieldnum"`
	Resulttype   Oid        `json:"resulttype"`
	Resulttypmod int32      `json:"resulttypmod"`
	Resultcollid Oid        `json:"resultcollid"`
}

func (n *FieldSelect) Pos() int {
	return 0
}
