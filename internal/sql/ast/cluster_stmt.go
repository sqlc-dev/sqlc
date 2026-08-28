package ast

type ClusterStmt struct {
	Tag NodeTag[ClusterStmt] `json:"tag"`

	Relation  *RangeVar
	Indexname *string
	Verbose   bool
}

func (n *ClusterStmt) Pos() int {
	return 0
}
