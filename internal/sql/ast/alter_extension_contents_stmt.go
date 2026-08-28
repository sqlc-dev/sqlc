package ast

type AlterExtensionContentsStmt struct {
	Tag NodeTag[AlterExtensionContentsStmt] `json:"tag"`

	Extname *string `json:",omitempty"`
	Action  int
	Objtype ObjectType
	Object  Node `json:",omitempty"`
}

func (n *AlterExtensionContentsStmt) Pos() int {
	return 0
}
