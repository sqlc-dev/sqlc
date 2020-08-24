package ast

type CreateEnumStmt_PG struct {
	TypeName *List
	Vals     *List
}

func (n *CreateEnumStmt_PG) Pos() int {
	return 0
}
