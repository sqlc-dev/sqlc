package ast

type CreateUserMappingStmt struct {
	Tag NodeTag[CreateUserMappingStmt] `json:"tag"`

	User        *RoleSpec `json:",omitempty"`
	Servername  *string   `json:",omitempty"`
	IfNotExists bool
	Options     *List `json:",omitempty"`
}

func (n *CreateUserMappingStmt) Pos() int {
	return 0
}
