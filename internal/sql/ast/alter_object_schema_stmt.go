package ast

type AlterObjectSchemaStmt struct {
	ObjectType ObjectType
	Relation   *RangeVar
	Object     Node
	Newschema  *string
	MissingOk  bool
}

func (n *AlterObjectSchemaStmt) Pos() int {
	return 0
}
