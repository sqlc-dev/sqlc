package ast

type ViewStmt struct {
	Tag NodeTag[ViewStmt] `json:"tag"`

	View            *RangeVar       `json:"view,omitempty"`
	Aliases         *List           `json:"aliases,omitempty"`
	Query           Node            `json:"query,omitempty"`
	Replace         bool            `json:"replace"`
	Options         *List           `json:"options,omitempty"`
	WithCheckOption ViewCheckOption `json:"with_check_option"`
}

func (n *ViewStmt) Pos() int {
	return 0
}
