package ast

type CreateCastStmt struct {
	Tag NodeTag[CreateCastStmt] `json:"tag"`

	Sourcetype *TypeName       `json:"sourcetype,omitempty"`
	Targettype *TypeName       `json:"targettype,omitempty"`
	Func       *ObjectWithArgs `json:"func,omitempty"`
	Context    CoercionContext `json:"context"`
	Inout      bool            `json:"inout"`
}

func (n *CreateCastStmt) Pos() int {
	return 0
}
