package ast

type AlterEnumStmt struct {
	Tag NodeTag[AlterEnumStmt] `json:"tag"`

	TypeName           *List
	OldVal             *string
	NewVal             *string
	NewValNeighbor     *string
	NewValIsAfter      bool
	SkipIfNewValExists bool
}

func (n *AlterEnumStmt) Pos() int {
	return 0
}
