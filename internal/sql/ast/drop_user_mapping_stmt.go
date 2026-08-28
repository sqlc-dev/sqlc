package ast

type DropUserMappingStmt struct {
	Tag NodeTag[DropUserMappingStmt] `json:"tag"`

	User       *RoleSpec `json:"user,omitempty"`
	Servername *string   `json:"servername,omitempty"`
	MissingOk  bool      `json:"missing_ok"`
}

func (n *DropUserMappingStmt) Pos() int {
	return 0
}
