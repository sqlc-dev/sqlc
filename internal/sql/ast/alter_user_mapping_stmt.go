package ast

type AlterUserMappingStmt struct {
	Tag NodeTag[AlterUserMappingStmt] `json:"tag"`

	User       *RoleSpec
	Servername *string
	Options    *List
}

func (n *AlterUserMappingStmt) Pos() int {
	return 0
}
