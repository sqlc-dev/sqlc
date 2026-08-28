package ast

type TableFunc struct {
	Tag NodeTag[TableFunc] `json:"tag"`

	NsUris        *List    `json:"ns_uris,omitempty"`
	NsNames       *List    `json:"ns_names,omitempty"`
	Docexpr       Node     `json:"docexpr,omitempty"`
	Rowexpr       Node     `json:"rowexpr,omitempty"`
	Colnames      *List    `json:"colnames,omitempty"`
	Coltypes      *List    `json:"coltypes,omitempty"`
	Coltypmods    *List    `json:"coltypmods,omitempty"`
	Colcollations *List    `json:"colcollations,omitempty"`
	Colexprs      *List    `json:"colexprs,omitempty"`
	Coldefexprs   *List    `json:"coldefexprs,omitempty"`
	Notnulls      []uint32 `json:"notnulls,omitempty"`
	Ordinalitycol int      `json:"ordinalitycol"`
	Location      int      `json:"location"`
}

func (n *TableFunc) Pos() int {
	return n.Location
}
