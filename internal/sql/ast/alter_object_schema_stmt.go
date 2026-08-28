package ast

type AlterObjectSchemaStmt struct {
	Tag NodeTag[AlterObjectSchemaStmt] `json:"tag"`

	ObjectType ObjectType `json:"object_type"`
	Relation   *RangeVar  `json:"relation,omitempty"`
	Object     Node       `json:"object,omitempty"`
	Newschema  *string    `json:"newschema,omitempty"`
	MissingOk  bool       `json:"missing_ok"`
}

func (n *AlterObjectSchemaStmt) Pos() int {
	return 0
}
