package ast

type SecLabelStmt struct {
	Tag NodeTag[SecLabelStmt] `json:"tag"`

	Objtype  ObjectType `json:"objtype"`
	Object   Node       `json:"object,omitempty"`
	Provider *string    `json:"provider,omitempty"`
	Label    *string    `json:"label,omitempty"`
}

func (n *SecLabelStmt) Pos() int {
	return 0
}
