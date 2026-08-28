package ast

type FuncSpec struct {
	Tag NodeTag[FuncSpec] `json:"tag"`

	Name    *FuncName
	Args    []*TypeName
	HasArgs bool
}

func (n *FuncSpec) Pos() int {
	return 0
}
