package ast

type CreatePolicyStmt struct {
	Tag NodeTag[CreatePolicyStmt] `json:"tag"`

	PolicyName *string   `json:",omitempty"`
	Table      *RangeVar `json:",omitempty"`
	CmdName    *string   `json:",omitempty"`
	Permissive bool
	Roles      *List `json:",omitempty"`
	Qual       Node  `json:",omitempty"`
	WithCheck  Node  `json:",omitempty"`
}

func (n *CreatePolicyStmt) Pos() int {
	return 0
}
