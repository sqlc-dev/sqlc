package ast

type AlterPolicyStmt struct {
	Tag NodeTag[AlterPolicyStmt] `json:"tag"`

	PolicyName *string
	Table      *RangeVar
	Roles      *List
	Qual       Node
	WithCheck  Node
}

func (n *AlterPolicyStmt) Pos() int {
	return 0
}
