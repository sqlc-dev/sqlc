package ast

type CreateOpFamilyStmt struct {
	Tag NodeTag[CreateOpFamilyStmt] `json:"tag"`

	Opfamilyname *List
	Amname       *string
}

func (n *CreateOpFamilyStmt) Pos() int {
	return 0
}
