package ast

type CreateUserMappingStmt struct {
	Tag NodeTag[CreateUserMappingStmt] `json:"tag"`

	User        *RoleSpec `json:"user,omitempty"`
	Servername  *string   `json:"servername,omitempty"`
	IfNotExists bool      `json:"if_not_exists"`
	Options     *List     `json:"options,omitempty"`
}

func (n *CreateUserMappingStmt) Pos() int {
	return 0
}
