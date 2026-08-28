package ast

type CreateOpClassStmt struct {
	Tag NodeTag[CreateOpClassStmt] `json:"tag"`

	Opclassname  *List     `json:"opclassname,omitempty"`
	Opfamilyname *List     `json:"opfamilyname,omitempty"`
	Amname       *string   `json:"amname,omitempty"`
	Datatype     *TypeName `json:"datatype,omitempty"`
	Items        *List     `json:"items,omitempty"`
	IsDefault    bool      `json:"is_default"`
}

func (n *CreateOpClassStmt) Pos() int {
	return 0
}
