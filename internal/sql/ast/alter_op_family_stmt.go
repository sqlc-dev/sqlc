package ast

type AlterOpFamilyStmt struct {
	Tag NodeTag[AlterOpFamilyStmt] `json:"tag"`

	Opfamilyname *List   `json:"opfamilyname,omitempty"`
	Amname       *string `json:"amname,omitempty"`
	IsDrop       bool    `json:"is_drop"`
	Items        *List   `json:"items,omitempty"`
}

func (n *AlterOpFamilyStmt) Pos() int {
	return 0
}
