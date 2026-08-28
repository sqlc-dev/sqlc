package ast

type AlterDefaultPrivilegesStmt struct {
	Tag NodeTag[AlterDefaultPrivilegesStmt] `json:"tag"`

	Options *List
	Action  *GrantStmt
}

func (n *AlterDefaultPrivilegesStmt) Pos() int {
	return 0
}
