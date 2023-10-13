package ast

type AlterTableStmt struct {
	// TODO: Only TableName or Relation should be defined
	Relation  *RangeVar
	Table     *TableName
	Cmds      *List
	MissingOk bool
	Relkind   ObjectType
}

func (n *AlterTableStmt) Pos() int {
	return 0
}

func (n *AlterTableStmt) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}
	buf.WriteString("ALTER TABLE ")
	buf.astFormat(n.Relation)
	buf.astFormat(n.Table)
	buf.astFormat(n.Cmds)
}
