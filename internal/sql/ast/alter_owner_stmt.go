package ast

type AlterOwnerStmt struct {
	Tag NodeTag[AlterOwnerStmt] `json:"tag"`

	ObjectType ObjectType
	Relation   *RangeVar `json:",omitempty"`
	Object     Node      `json:",omitempty"`
	Newowner   *RoleSpec `json:",omitempty"`
}

func (n *AlterOwnerStmt) Pos() int {
	return 0
}
