package ast

type MatViewStmt struct {
	MatView         *RangeVar
	Aliases         *List
	Query           Node
	Replace         bool
	Options         *List
	WithCheckOption ViewCheckOption
}

func (n *MatViewStmt) Pos() int {
	return 0
}
