package ast

type CreateEnumStmt struct {
	Tag NodeTag[CreateEnumStmt] `json:"tag"`

	TypeName *TypeName
	Vals     *List
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
