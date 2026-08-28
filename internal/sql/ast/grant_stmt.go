package ast

type GrantStmt struct {
	Tag NodeTag[GrantStmt] `json:"tag"`

	IsGrant     bool
	Targtype    GrantTargetType
	Objtype     GrantObjectType
	Objects     *List `json:",omitempty"`
	Privileges  *List `json:",omitempty"`
	Grantees    *List `json:",omitempty"`
	GrantOption bool
	Behavior    DropBehavior
}

func (n *GrantStmt) Pos() int {
	return 0
}
