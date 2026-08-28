package ast

type CreateTableSpaceStmt struct {
	Tag NodeTag[CreateTableSpaceStmt] `json:"tag"`

	Tablespacename *string
	Owner          *RoleSpec
	Location       *string
	Options        *List
}

func (n *CreateTableSpaceStmt) Pos() int {
	return 0
}
