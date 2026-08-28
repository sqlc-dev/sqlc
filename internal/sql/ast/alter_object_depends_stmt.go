package ast

type AlterObjectDependsStmt struct {
	Tag NodeTag[AlterObjectDependsStmt] `json:"tag"`

	ObjectType ObjectType
	Relation   *RangeVar `json:",omitempty"`
	Object     Node      `json:",omitempty"`
	Extname    Node      `json:",omitempty"`
}

func (n *AlterObjectDependsStmt) Pos() int {
	return 0
}
