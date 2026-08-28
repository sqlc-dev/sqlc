package ast

type CreateUserMappingStmt struct {
	Tag NodeTag[CreateUserMappingStmt] `json:"tag"`

	User        *RoleSpec
	Servername  *string
	IfNotExists bool
	Options     *List
}

func (n *CreateUserMappingStmt) Pos() int {
	return 0
}
