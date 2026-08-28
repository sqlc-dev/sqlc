package ast

type CoerceToDomain struct {
	Tag NodeTag[CoerceToDomain] `json:"tag"`

	Xpr            Node         `json:"xpr,omitempty"`
	Arg            Node         `json:"arg,omitempty"`
	Resulttype     Oid          `json:"resulttype"`
	Resulttypmod   int32        `json:"resulttypmod"`
	Resultcollid   Oid          `json:"resultcollid"`
	Coercionformat CoercionForm `json:"coercionformat"`
	Location       int          `json:"location"`
}

func (n *CoerceToDomain) Pos() int {
	return n.Location
}
