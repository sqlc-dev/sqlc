package ast

type GrantStmt struct {
	IsGrant     bool
	Targtype    GrantTargetType
	Objtype     GrantObjectType
	Objects     *List
	Privileges  *List
	Grantees    *List
	GrantOption bool
	Behavior    DropBehavior
}

func (n *GrantStmt) Pos() int {
	return 0
}
