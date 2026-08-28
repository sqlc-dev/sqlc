package ast

type RangeTblFunction struct {
	Tag NodeTag[RangeTblFunction] `json:"tag"`

	Funcexpr          Node `json:",omitempty"`
	Funccolcount      int
	Funccolnames      *List    `json:",omitempty"`
	Funccoltypes      *List    `json:",omitempty"`
	Funccoltypmods    *List    `json:",omitempty"`
	Funccolcollations *List    `json:",omitempty"`
	Funcparams        []uint32 `json:",omitempty"`
}

func (n *RangeTblFunction) Pos() int {
	return 0
}
