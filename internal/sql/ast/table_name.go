package ast

type TableName struct {
	Catalog string
	Schema  string
	Name    string
}

func (n *TableName) Pos() int {
	return 0
}
