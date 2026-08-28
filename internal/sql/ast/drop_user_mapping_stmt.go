package ast

type DropUserMappingStmt struct {
	Tag NodeTag[DropUserMappingStmt] `json:"tag"`

	User       *RoleSpec
	Servername *string
	MissingOk  bool
}

func (n *DropUserMappingStmt) Pos() int {
	return 0
}
