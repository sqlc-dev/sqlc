package ast

type TODO struct {
	Tag NodeTag[TODO] `json:"tag"`
}

func (n *TODO) Pos() int {
	return 0
}
