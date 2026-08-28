package ast

type AlterPolicyStmt struct {
	Tag NodeTag[AlterPolicyStmt] `json:"tag"`

	PolicyName *string   `json:",omitempty"`
	Table      *RangeVar `json:",omitempty"`
	Roles      *List     `json:",omitempty"`
	Qual       Node      `json:",omitempty"`
	WithCheck  Node      `json:",omitempty"`
}

func (n *AlterPolicyStmt) Pos() int {
	return 0
}
