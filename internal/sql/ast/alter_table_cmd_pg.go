package ast

type AlterTableCmd_PG struct {
	Subtype   AlterTableType
	Name      *string
	Newowner  *RoleSpec
	Def       Node
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterTableCmd_PG) Pos() int {
	return 0
}
