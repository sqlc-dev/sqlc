package ast

type AlterUserMappingStmt struct {
	Tag NodeTag[AlterUserMappingStmt] `json:"tag"`

	User       *RoleSpec `json:"user,omitempty"`
	Servername *string   `json:"servername,omitempty"`
	Options    *List     `json:"options,omitempty"`
}

func (n *AlterUserMappingStmt) Pos() int {
	return 0
}
