package ast

type CurrentOfExpr struct {
	Tag NodeTag[CurrentOfExpr] `json:"tag"`

	Xpr         Node
	Cvarno      Index
	CursorName  *string
	CursorParam int
}

func (n *CurrentOfExpr) Pos() int {
	return 0
}
