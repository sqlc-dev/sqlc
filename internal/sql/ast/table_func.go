package ast

type TableFunc struct {
	NsUris        *List
	NsNames       *List
	Docexpr       Node
	Rowexpr       Node
	Colnames      *List
	Coltypes      *List
	Coltypmods    *List
	Colcollations *List
	Colexprs      *List
	Coldefexprs   *List
	Notnulls      []uint32
	Ordinalitycol int
	Location      int
}

func (n *TableFunc) Pos() int {
	return n.Location
}
