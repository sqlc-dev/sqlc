package ast

type CreateEnumStmt struct {
	Tag NodeTag[CreateEnumStmt] `json:"tag"`

	TypeName *TypeName `json:"type_name,omitempty"`
	Vals     *List     `json:"vals,omitempty"`
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
