package ast

type DropUserMappingStmt struct {
	Tag NodeTag[DropUserMappingStmt] `json:"tag"`

	User       *RoleSpec `json:",omitempty"`
	Servername *string   `json:",omitempty"`
	MissingOk  bool
}

func (n *DropUserMappingStmt) Pos() int {
	return 0
}
