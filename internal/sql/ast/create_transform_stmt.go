package ast

type CreateTransformStmt struct {
	Tag NodeTag[CreateTransformStmt] `json:"tag"`

	Replace  bool            `json:"replace"`
	TypeName *TypeName       `json:"type_name,omitempty"`
	Lang     *string         `json:"lang,omitempty"`
	Fromsql  *ObjectWithArgs `json:"fromsql,omitempty"`
	Tosql    *ObjectWithArgs `json:"tosql,omitempty"`
}

func (n *CreateTransformStmt) Pos() int {
	return 0
}
