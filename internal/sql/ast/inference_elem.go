package ast

type InferenceElem struct {
	Tag NodeTag[InferenceElem] `json:"tag"`

	Xpr          Node `json:",omitempty"`
	Expr         Node `json:",omitempty"`
	Infercollid  Oid
	Inferopclass Oid
}

func (n *InferenceElem) Pos() int {
	return 0
}
