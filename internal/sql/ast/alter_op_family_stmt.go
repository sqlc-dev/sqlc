package ast

type AlterOpFamilyStmt struct {
	Tag NodeTag[AlterOpFamilyStmt] `json:"tag"`

	Opfamilyname *List
	Amname       *string
	IsDrop       bool
	Items        *List
}

func (n *AlterOpFamilyStmt) Pos() int {
	return 0
}
