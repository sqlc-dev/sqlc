package ast

type CoerceToDomain struct {
	Tag NodeTag[CoerceToDomain] `json:"tag"`

	Xpr            Node
	Arg            Node
	Resulttype     Oid
	Resulttypmod   int32
	Resultcollid   Oid
	Coercionformat CoercionForm
	Location       int
}

func (n *CoerceToDomain) Pos() int {
	return n.Location
}
