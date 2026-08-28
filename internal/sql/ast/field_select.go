package ast

type FieldSelect struct {
	Tag NodeTag[FieldSelect] `json:"tag"`

	Xpr          Node
	Arg          Node
	Fieldnum     AttrNumber
	Resulttype   Oid
	Resulttypmod int32
	Resultcollid Oid
}

func (n *FieldSelect) Pos() int {
	return 0
}
