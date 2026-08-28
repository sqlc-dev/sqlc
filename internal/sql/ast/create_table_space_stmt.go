package ast

type CreateTableSpaceStmt struct {
	Tag NodeTag[CreateTableSpaceStmt] `json:"tag"`

	Tablespacename *string   `json:"tablespacename,omitempty"`
	Owner          *RoleSpec `json:"owner,omitempty"`
	Location       *string   `json:"location,omitempty"`
	Options        *List     `json:"options,omitempty"`
}

func (n *CreateTableSpaceStmt) Pos() int {
	return 0
}
