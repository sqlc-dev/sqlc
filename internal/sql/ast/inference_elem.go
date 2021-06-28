package ast

type InferenceElem struct {
	Xpr          Node
	Expr         Node
	Infercollid  Oid
	Inferopclass Oid
}

func (n *InferenceElem) Pos() int {
	return 0
}
