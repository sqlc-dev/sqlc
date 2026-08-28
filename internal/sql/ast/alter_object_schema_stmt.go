package ast

type AlterObjectSchemaStmt struct {
	Tag NodeTag[AlterObjectSchemaStmt] `json:"tag"`

	ObjectType ObjectType
	Relation   *RangeVar `json:",omitempty"`
	Object     Node      `json:",omitempty"`
	Newschema  *string   `json:",omitempty"`
	MissingOk  bool
}

func (n *AlterObjectSchemaStmt) Pos() int {
	return 0
}
