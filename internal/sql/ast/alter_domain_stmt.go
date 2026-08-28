package ast

type AlterDomainStmt struct {
	Tag NodeTag[AlterDomainStmt] `json:"tag"`

	Subtype   byte
	TypeName  *List   `json:",omitempty"`
	Name      *string `json:",omitempty"`
	Def       Node    `json:",omitempty"`
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterDomainStmt) Pos() int {
	return 0
}
