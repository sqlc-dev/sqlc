package ast

type RangeTblFunction struct {
	Tag NodeTag[RangeTblFunction] `json:"tag"`

	Funcexpr          Node
	Funccolcount      int
	Funccolnames      *List
	Funccoltypes      *List
	Funccoltypmods    *List
	Funccolcollations *List
	Funcparams        []uint32
}

func (n *RangeTblFunction) Pos() int {
	return 0
}
