package ast

type AlterObjectDependsStmt struct {
	Tag NodeTag[AlterObjectDependsStmt] `json:"tag"`

	ObjectType ObjectType `json:"object_type"`
	Relation   *RangeVar  `json:"relation,omitempty"`
	Object     Node       `json:"object,omitempty"`
	Extname    Node       `json:"extname,omitempty"`
}

func (n *AlterObjectDependsStmt) Pos() int {
	return 0
}
