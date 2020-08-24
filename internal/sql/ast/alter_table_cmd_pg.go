package ast

import ()

type AlterTableCmd struct {
	Subtype   AlterTableType
	Name      *string
	Newowner  *RoleSpec
	Def       Node
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterTableCmd) Pos() int {
	return 0
}
