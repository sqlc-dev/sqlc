package ast

type FuncSpec struct {
	Tag NodeTag[FuncSpec] `json:"tag"`

	Name    *FuncName   `json:",omitempty"`
	Args    []*TypeName `json:",omitempty"`
	HasArgs bool
}

func (n *FuncSpec) Pos() int {
	return 0
}
