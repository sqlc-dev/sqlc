package ast

type InferenceElem struct {
	Tag NodeTag[InferenceElem] `json:"tag"`

	Xpr          Node
	Expr         Node
	Infercollid  Oid
	Inferopclass Oid
}

func (n *InferenceElem) Pos() int {
	return 0
}
