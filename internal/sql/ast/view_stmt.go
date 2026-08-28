package ast

type ViewStmt struct {
	Tag NodeTag[ViewStmt] `json:"tag"`

	View            *RangeVar `json:",omitempty"`
	Aliases         *List     `json:",omitempty"`
	Query           Node      `json:",omitempty"`
	Replace         bool
	Options         *List `json:",omitempty"`
	WithCheckOption ViewCheckOption
}

func (n *ViewStmt) Pos() int {
	return 0
}
