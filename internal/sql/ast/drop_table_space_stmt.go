package ast

type DropTableSpaceStmt struct {
	Tag NodeTag[DropTableSpaceStmt] `json:"tag"`

	Tablespacename *string `json:",omitempty"`
	MissingOk      bool
}

func (n *DropTableSpaceStmt) Pos() int {
	return 0
}
