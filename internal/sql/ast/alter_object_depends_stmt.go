package ast

type AlterObjectDependsStmt struct {
	Tag NodeTag[AlterObjectDependsStmt] `json:"tag"`

	ObjectType ObjectType
	Relation   *RangeVar
	Object     Node
	Extname    Node
}

func (n *AlterObjectDependsStmt) Pos() int {
	return 0
}
