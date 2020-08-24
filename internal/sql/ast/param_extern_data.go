package ast

type ParamExternData struct {
	Value  Datum
	Isnull bool
	Pflags uint16
	Ptype  Oid
}

func (n *ParamExternData) Pos() int {
	return 0
}
