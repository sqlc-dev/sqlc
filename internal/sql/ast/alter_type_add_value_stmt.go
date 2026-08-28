package ast

type AlterTypeAddValueStmt struct {
	Tag NodeTag[AlterTypeAddValueStmt] `json:"tag"`

	Type               *TypeName
	NewValue           *string
	NewValHasNeighbor  bool
	NewValNeighbor     *string
	NewValIsAfter      bool
	SkipIfNewValExists bool
}

func (n *AlterTypeAddValueStmt) Pos() int {
	return 0
}
