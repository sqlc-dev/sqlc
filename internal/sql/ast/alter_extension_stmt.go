package ast

type AlterExtensionStmt struct {
	Tag NodeTag[AlterExtensionStmt] `json:"tag"`

	Extname *string `json:"extname,omitempty"`
	Options *List   `json:"options,omitempty"`
}

func (n *AlterExtensionStmt) Pos() int {
	return 0
}
