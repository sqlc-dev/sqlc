package ast

type AlterDefaultPrivilegesStmt struct {
	Tag NodeTag[AlterDefaultPrivilegesStmt] `json:"tag"`

	Options *List      `json:"options,omitempty"`
	Action  *GrantStmt `json:"action,omitempty"`
}

func (n *AlterDefaultPrivilegesStmt) Pos() int {
	return 0
}
