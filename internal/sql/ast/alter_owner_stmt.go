package ast

type AlterOwnerStmt struct {
	ObjectType ObjectType
	Relation   *RangeVar
	Object     Node
	Newowner   *RoleSpec
}

func (n *AlterOwnerStmt) Pos() int {
	return 0
}
