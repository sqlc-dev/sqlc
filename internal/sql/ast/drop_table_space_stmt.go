package ast

type DropTableSpaceStmt struct {
	Tag NodeTag[DropTableSpaceStmt] `json:"tag"`

	Tablespacename *string `json:"tablespacename,omitempty"`
	MissingOk      bool    `json:"missing_ok"`
}

func (n *DropTableSpaceStmt) Pos() int {
	return 0
}
