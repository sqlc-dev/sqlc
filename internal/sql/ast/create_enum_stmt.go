package ast

type CreateEnumStmt struct {
	Tag NodeTag[CreateEnumStmt] `json:"tag"`

	TypeName *TypeName `json:",omitempty"`
	Vals     *List     `json:",omitempty"`
}

func (n *CreateEnumStmt) Pos() int {
	return 0
}
