package ast

type AlterOpFamilyStmt struct {
	Tag NodeTag[AlterOpFamilyStmt] `json:"tag"`

	Opfamilyname *List   `json:",omitempty"`
	Amname       *string `json:",omitempty"`
	IsDrop       bool
	Items        *List `json:",omitempty"`
}

func (n *AlterOpFamilyStmt) Pos() int {
	return 0
}
