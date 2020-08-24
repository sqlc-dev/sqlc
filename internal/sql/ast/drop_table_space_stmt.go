package ast

type DropTableSpaceStmt struct {
	Tablespacename *string
	MissingOk      bool
}

func (n *DropTableSpaceStmt) Pos() int {
	return 0
}
