package ast

type BitString struct {
	Tag NodeTag[BitString] `json:"tag"`

	Str string
}

func (n *BitString) Pos() int {
	return 0
}
