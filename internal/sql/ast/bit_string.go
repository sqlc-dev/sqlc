package ast

type BitString struct {
	Tag NodeTag[BitString] `json:"tag"`

	Str string `json:"str"`
}

func (n *BitString) Pos() int {
	return 0
}
