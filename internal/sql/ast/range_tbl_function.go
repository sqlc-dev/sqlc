package ast

type RangeTblFunction struct {
	Tag NodeTag[RangeTblFunction] `json:"tag"`

	Funcexpr          Node     `json:"funcexpr,omitempty"`
	Funccolcount      int      `json:"funccolcount"`
	Funccolnames      *List    `json:"funccolnames,omitempty"`
	Funccoltypes      *List    `json:"funccoltypes,omitempty"`
	Funccoltypmods    *List    `json:"funccoltypmods,omitempty"`
	Funccolcollations *List    `json:"funccolcollations,omitempty"`
	Funcparams        []uint32 `json:"funcparams,omitempty"`
}

func (n *RangeTblFunction) Pos() int {
	return 0
}
