package ast

type AlterCollationStmt struct {
	Tag NodeTag[AlterCollationStmt] `json:"tag"`

	Collname *List `json:"collname,omitempty"`
}

func (n *AlterCollationStmt) Pos() int {
	return 0
}
