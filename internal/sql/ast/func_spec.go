package ast

type FuncSpec struct {
	Name    *FuncName
	Args    []*TypeName
	HasArgs bool
}

func (n *FuncSpec) Pos() int {
	return 0
}
