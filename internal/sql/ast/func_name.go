package ast

type FuncName struct {
	Catalog string
	Schema  string
	Name    string
}

func (n *FuncName) Pos() int {
	return 0
}
