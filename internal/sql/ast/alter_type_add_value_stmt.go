package ast

type AlterTypeAddValueStmt struct {
	Tag NodeTag[AlterTypeAddValueStmt] `json:"tag"`

	Type               *TypeName `json:",omitempty"`
	NewValue           *string   `json:",omitempty"`
	NewValHasNeighbor  bool
	NewValNeighbor     *string `json:",omitempty"`
	NewValIsAfter      bool
	SkipIfNewValExists bool
}

func (n *AlterTypeAddValueStmt) Pos() int {
	return 0
}
