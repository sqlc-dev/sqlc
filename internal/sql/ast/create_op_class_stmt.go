package ast

type CreateOpClassStmt struct {
	Tag NodeTag[CreateOpClassStmt] `json:"tag"`

	Opclassname  *List     `json:",omitempty"`
	Opfamilyname *List     `json:",omitempty"`
	Amname       *string   `json:",omitempty"`
	Datatype     *TypeName `json:",omitempty"`
	Items        *List     `json:",omitempty"`
	IsDefault    bool
}

func (n *CreateOpClassStmt) Pos() int {
	return 0
}
