package ast

type ViewStmt struct {
	Tag NodeTag[ViewStmt] `json:"tag"`

	View            *RangeVar
	Aliases         *List
	Query           Node
	Replace         bool
	Options         *List
	WithCheckOption ViewCheckOption
}

func (n *ViewStmt) Pos() int {
	return 0
}
