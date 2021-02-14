package ast

type SecLabelStmt struct {
	Objtype  ObjectType
	Object   Node
	Provider *string
	Label    *string
}

func (n *SecLabelStmt) Pos() int {
	return 0
}
