package ast

type AlterUserMappingStmt struct {
	User       *RoleSpec
	Servername *string
	Options    *List
}

func (n *AlterUserMappingStmt) Pos() int {
	return 0
}
