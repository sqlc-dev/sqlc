package ast

type ParamExternData struct {
	Tag NodeTag[ParamExternData] `json:"tag"`

	Value  Datum  `json:"value"`
	Isnull bool   `json:"isnull"`
	Pflags uint16 `json:"pflags"`
	Ptype  Oid    `json:"ptype"`
}

func (n *ParamExternData) Pos() int {
	return 0
}
