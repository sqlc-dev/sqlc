package pg

type DropUserMappingStmt struct {
	User       *RoleSpec
	Servername *string
	MissingOk  bool
}

func (n *DropUserMappingStmt) Pos() int {
	return 0
}
