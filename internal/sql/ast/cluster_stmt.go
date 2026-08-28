package ast

type ClusterStmt struct {
	Tag NodeTag[ClusterStmt] `json:"tag"`

	Relation  *RangeVar `json:",omitempty"`
	Indexname *string   `json:",omitempty"`
	Verbose   bool
}

func (n *ClusterStmt) Pos() int {
	return 0
}
