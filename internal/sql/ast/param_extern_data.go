package ast

type ParamExternData struct {
	Tag NodeTag[ParamExternData] `json:"tag"`

	Value  Datum
	Isnull bool
	Pflags uint16
	Ptype  Oid
}

func (n *ParamExternData) Pos() int {
	return 0
}
