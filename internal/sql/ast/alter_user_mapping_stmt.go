package ast

type AlterUserMappingStmt struct {
	Tag NodeTag[AlterUserMappingStmt] `json:"tag"`

	User       *RoleSpec `json:",omitempty"`
	Servername *string   `json:",omitempty"`
	Options    *List     `json:",omitempty"`
}

func (n *AlterUserMappingStmt) Pos() int {
	return 0
}
