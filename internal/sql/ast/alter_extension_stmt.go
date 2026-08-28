package ast

type AlterExtensionStmt struct {
	Tag NodeTag[AlterExtensionStmt] `json:"tag"`

	Extname *string
	Options *List
}

func (n *AlterExtensionStmt) Pos() int {
	return 0
}
