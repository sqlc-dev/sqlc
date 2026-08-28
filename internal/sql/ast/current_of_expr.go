package ast

type CurrentOfExpr struct {
	Tag NodeTag[CurrentOfExpr] `json:"tag"`

	Xpr         Node    `json:"xpr,omitempty"`
	Cvarno      Index   `json:"cvarno"`
	CursorName  *string `json:"cursor_name,omitempty"`
	CursorParam int     `json:"cursor_param"`
}

func (n *CurrentOfExpr) Pos() int {
	return 0
}
