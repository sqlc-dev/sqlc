package ast

type CreateEnumStmt struct {
	TypeName *TypeName
	Vals     *List
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
