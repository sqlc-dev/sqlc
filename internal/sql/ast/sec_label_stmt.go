package ast

type SecLabelStmt struct {
	Tag NodeTag[SecLabelStmt] `json:"tag"`

	Objtype  ObjectType
	Object   Node    `json:",omitempty"`
	Provider *string `json:",omitempty"`
	Label    *string `json:",omitempty"`
}

func (n *SecLabelStmt) Pos() int {
	return 0
}
