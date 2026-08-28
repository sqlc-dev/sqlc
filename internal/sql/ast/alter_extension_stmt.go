package ast

type AlterExtensionStmt struct {
	Tag NodeTag[AlterExtensionStmt] `json:"tag"`

	Extname *string `json:",omitempty"`
	Options *List   `json:",omitempty"`
}

func (n *AlterExtensionStmt) Pos() int {
	return 0
}
