package ast

type CurrentOfExpr struct {
	Tag NodeTag[CurrentOfExpr] `json:"tag"`

	Xpr         Node `json:",omitempty"`
	Cvarno      Index
	CursorName  *string `json:",omitempty"`
	CursorParam int
}

func (n *CurrentOfExpr) Pos() int {
	return 0
}
