package ast

type AlterDomainStmt struct {
	Tag NodeTag[AlterDomainStmt] `json:"tag"`

	Subtype   byte         `json:"subtype"`
	TypeName  *List        `json:"type_name,omitempty"`
	Name      *string      `json:"name,omitempty"`
	Def       Node         `json:"def,omitempty"`
	Behavior  DropBehavior `json:"behavior"`
	MissingOk bool         `json:"missing_ok"`
}

func (n *AlterDomainStmt) Pos() int {
	return 0
}
