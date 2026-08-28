package ast

type AlterPolicyStmt struct {
	Tag NodeTag[AlterPolicyStmt] `json:"tag"`

	PolicyName *string   `json:"policy_name,omitempty"`
	Table      *RangeVar `json:"table,omitempty"`
	Roles      *List     `json:"roles,omitempty"`
	Qual       Node      `json:"qual,omitempty"`
	WithCheck  Node      `json:"with_check,omitempty"`
}

func (n *AlterPolicyStmt) Pos() int {
	return 0
}
