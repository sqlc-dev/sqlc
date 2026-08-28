package ast

type InferenceElem struct {
	Tag NodeTag[InferenceElem] `json:"tag"`

	Xpr          Node `json:"xpr,omitempty"`
	Expr         Node `json:"expr,omitempty"`
	Infercollid  Oid  `json:"infercollid"`
	Inferopclass Oid  `json:"inferopclass"`
}

func (n *InferenceElem) Pos() int {
	return 0
}
