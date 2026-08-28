package ast

type AlterExtensionContentsStmt struct {
	Tag NodeTag[AlterExtensionContentsStmt] `json:"tag"`

	Extname *string    `json:"extname,omitempty"`
	Action  int        `json:"action"`
	Objtype ObjectType `json:"objtype"`
	Object  Node       `json:"object,omitempty"`
}

func (n *AlterExtensionContentsStmt) Pos() int {
	return 0
}
