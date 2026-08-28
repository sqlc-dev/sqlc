package ast

type TableFunc struct {
	Tag NodeTag[TableFunc] `json:"tag"`

	NsUris        *List    `json:",omitempty"`
	NsNames       *List    `json:",omitempty"`
	Docexpr       Node     `json:",omitempty"`
	Rowexpr       Node     `json:",omitempty"`
	Colnames      *List    `json:",omitempty"`
	Coltypes      *List    `json:",omitempty"`
	Coltypmods    *List    `json:",omitempty"`
	Colcollations *List    `json:",omitempty"`
	Colexprs      *List    `json:",omitempty"`
	Coldefexprs   *List    `json:",omitempty"`
	Notnulls      []uint32 `json:",omitempty"`
	Ordinalitycol int
	Location      int
}

func (n *TableFunc) Pos() int {
	return n.Location
}
