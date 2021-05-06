package ast

type AlterExtensionContentsStmt struct {
	Extname *string
	Action  int
	Objtype ObjectType
	Object  Node
}

func (n *AlterExtensionContentsStmt) Pos() int {
	return 0
}
