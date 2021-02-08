package ast

type CreatePolicyStmt struct {
	PolicyName *string
	Table      *RangeVar
	CmdName    *string
	Permissive bool
	Roles      *List
	Qual       Node
	WithCheck  Node
}

func (n *CreatePolicyStmt) Pos() int {
	return 0
}
