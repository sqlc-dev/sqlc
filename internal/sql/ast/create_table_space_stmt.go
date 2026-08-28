package ast

type CreateTableSpaceStmt struct {
	Tag NodeTag[CreateTableSpaceStmt] `json:"tag"`

	Tablespacename *string   `json:",omitempty"`
	Owner          *RoleSpec `json:",omitempty"`
	Location       *string   `json:",omitempty"`
	Options        *List     `json:",omitempty"`
}

func (n *CreateTableSpaceStmt) Pos() int {
	return 0
}
