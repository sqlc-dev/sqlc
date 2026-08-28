package ast

type FuncSpec struct {
	Tag NodeTag[FuncSpec] `json:"tag"`

	Name    *FuncName   `json:"name,omitempty"`
	Args    []*TypeName `json:"args,omitempty"`
	HasArgs bool        `json:"has_args"`
}

func (n *FuncSpec) Pos() int {
	return 0
}
