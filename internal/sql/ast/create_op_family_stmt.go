package ast

type CreateOpFamilyStmt struct {
	Tag NodeTag[CreateOpFamilyStmt] `json:"tag"`

	Opfamilyname *List   `json:"opfamilyname,omitempty"`
	Amname       *string `json:"amname,omitempty"`
}

func (n *CreateOpFamilyStmt) Pos() int {
	return 0
}
