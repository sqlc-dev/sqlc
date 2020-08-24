package ast

type AlterTableType int

const (
	AT_AddColumn AlterTableType = iota
	AT_AlterColumnType
	AT_DropColumn
	AT_DropNotNull
	AT_SetNotNull
)

type AlterTableCmd struct {
	Subtype   AlterTableType
	Name      *string
	Def       Node
	Newowner  *RoleSpec
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterTableCmd) Pos() int {
	return 0
}
