package ast

type AlterExtensionContentsStmt struct {
	Tag NodeTag[AlterExtensionContentsStmt] `json:"tag"`

	Extname *string
	Action  int
	Objtype ObjectType
	Object  Node
}

func (n *AlterExtensionContentsStmt) Pos() int {
	return 0
}
