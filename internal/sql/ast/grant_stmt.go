package ast

type GrantStmt struct {
	Tag NodeTag[GrantStmt] `json:"tag"`

	IsGrant     bool            `json:"is_grant"`
	Targtype    GrantTargetType `json:"targtype"`
	Objtype     GrantObjectType `json:"objtype"`
	Objects     *List           `json:"objects,omitempty"`
	Privileges  *List           `json:"privileges,omitempty"`
	Grantees    *List           `json:"grantees,omitempty"`
	GrantOption bool            `json:"grant_option"`
	Behavior    DropBehavior    `json:"behavior"`
}

func (n *GrantStmt) Pos() int {
	return 0
}
