package ast

type CreateOpClassStmt struct {
	Tag NodeTag[CreateOpClassStmt] `json:"tag"`

	Opclassname  *List
	Opfamilyname *List
	Amname       *string
	Datatype     *TypeName
	Items        *List
	IsDefault    bool
}

func (n *CreateOpClassStmt) Pos() int {
	return 0
}
