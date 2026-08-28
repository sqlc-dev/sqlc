package ast

type CreateOpFamilyStmt struct {
	Tag NodeTag[CreateOpFamilyStmt] `json:"tag"`

	Opfamilyname *List   `json:",omitempty"`
	Amname       *string `json:",omitempty"`
}

func (n *CreateOpFamilyStmt) Pos() int {
	return 0
}
