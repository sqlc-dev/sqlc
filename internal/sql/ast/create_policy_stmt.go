package ast

type CreatePolicyStmt struct {
	Tag NodeTag[CreatePolicyStmt] `json:"tag"`

	PolicyName *string   `json:"policy_name,omitempty"`
	Table      *RangeVar `json:"table,omitempty"`
	CmdName    *string   `json:"cmd_name,omitempty"`
	Permissive bool      `json:"permissive"`
	Roles      *List     `json:"roles,omitempty"`
	Qual       Node      `json:"qual,omitempty"`
	WithCheck  Node      `json:"with_check,omitempty"`
}

func (n *CreatePolicyStmt) Pos() int {
	return 0
}
