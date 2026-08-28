package ast

type CreateCastStmt struct {
	Tag NodeTag[CreateCastStmt] `json:"tag"`

	Sourcetype *TypeName
	Targettype *TypeName
	Func       *ObjectWithArgs
	Context    CoercionContext
	Inout      bool
}

func (n *CreateCastStmt) Pos() int {
	return 0
}
