package ast

type CreateTransformStmt struct {
	Tag NodeTag[CreateTransformStmt] `json:"tag"`

	Replace  bool
	TypeName *TypeName       `json:",omitempty"`
	Lang     *string         `json:",omitempty"`
	Fromsql  *ObjectWithArgs `json:",omitempty"`
	Tosql    *ObjectWithArgs `json:",omitempty"`
}

func (n *CreateTransformStmt) Pos() int {
	return 0
}
