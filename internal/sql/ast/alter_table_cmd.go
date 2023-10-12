package ast

const (
	AT_AddColumn AlterTableType = iota
	AT_AlterColumnType
	AT_DropColumn
	AT_DropNotNull
	AT_SetNotNull
)

type AlterTableType int

func (t AlterTableType) String() string {
	switch t {
	case AT_AddColumn:
		return "AddColumn"
	case AT_AlterColumnType:
		return "AlterColumnType"
	case AT_DropColumn:
		return "DropColumn"
	case AT_DropNotNull:
		return "DropNotNull"
	case AT_SetNotNull:
		return "SetNotNull"
	default:
		return "Unknown"
	}
}

type AlterTableCmd struct {
	Subtype   AlterTableType
	Name      *string
	Def       *ColumnDef
	Newowner  *RoleSpec
	Behavior  DropBehavior
	MissingOk bool
}

func (n *AlterTableCmd) Pos() int {
	return 0
}

func (n *AlterTableCmd) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	switch n.Subtype {
	case AT_AddColumn:
		buf.WriteString(" ADD COLUMN ")
	case AT_DropColumn:
		buf.WriteString(" DROP COLUMN ")
	}

	buf.astFormat(n.Def)
}
