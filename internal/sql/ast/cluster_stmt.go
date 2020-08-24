package ast

type ClusterStmt struct {
	Relation  *RangeVar
	Indexname *string
	Verbose   bool
}

func (n *ClusterStmt) Pos() int {
	return 0
}
