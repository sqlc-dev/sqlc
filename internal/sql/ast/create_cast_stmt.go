package ast

type CreateCastStmt struct {
	Tag NodeTag[CreateCastStmt] `json:"tag"`

	Sourcetype *TypeName       `json:",omitempty"`
	Targettype *TypeName       `json:",omitempty"`
	Func       *ObjectWithArgs `json:",omitempty"`
	Context    CoercionContext
	Inout      bool
}

func (n *CreateCastStmt) Pos() int {
	return 0
}
