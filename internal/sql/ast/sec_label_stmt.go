package ast

type SecLabelStmt struct {
	Tag NodeTag[SecLabelStmt] `json:"tag"`

	Objtype  ObjectType
	Object   Node
	Provider *string
	Label    *string
}

func (n *SecLabelStmt) Pos() int {
	return 0
}
