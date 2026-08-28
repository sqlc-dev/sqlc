package ast

type CreateTransformStmt struct {
	Tag NodeTag[CreateTransformStmt] `json:"tag"`

	Replace  bool
	TypeName *TypeName
	Lang     *string
	Fromsql  *ObjectWithArgs
	Tosql    *ObjectWithArgs
}

func (n *CreateTransformStmt) Pos() int {
	return 0
}
