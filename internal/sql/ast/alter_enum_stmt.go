package ast

type AlterEnumStmt struct {
	Tag NodeTag[AlterEnumStmt] `json:"tag"`

	TypeName           *List   `json:",omitempty"`
	OldVal             *string `json:",omitempty"`
	NewVal             *string `json:",omitempty"`
	NewValNeighbor     *string `json:",omitempty"`
	NewValIsAfter      bool
	SkipIfNewValExists bool
}

func (n *AlterEnumStmt) Pos() int {
	return 0
}
