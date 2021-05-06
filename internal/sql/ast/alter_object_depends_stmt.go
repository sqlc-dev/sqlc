package ast

type AlterObjectDependsStmt struct {
	ObjectType ObjectType
	Relation   *RangeVar
	Object     Node
	Extname    Node
}

func (n *AlterObjectDependsStmt) Pos() int {
	return 0
}
