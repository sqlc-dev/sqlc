package ast

type AlterOwnerStmt struct {
	Tag NodeTag[AlterOwnerStmt] `json:"tag"`

	ObjectType ObjectType
	Relation   *RangeVar
	Object     Node
	Newowner   *RoleSpec
}

func (n *AlterOwnerStmt) Pos() int {
	return 0
}
