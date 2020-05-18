package ast

type TypeName struct {
	Catalog string
	Schema  string
	Name    string
}

func (n *TypeName) Pos() int {
	return 0
}
