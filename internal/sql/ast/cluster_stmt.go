package ast

type ClusterStmt struct {
	Tag NodeTag[ClusterStmt] `json:"tag"`

	Relation  *RangeVar `json:"relation,omitempty"`
	Indexname *string   `json:"indexname,omitempty"`
	Verbose   bool      `json:"verbose"`
}

func (n *ClusterStmt) Pos() int {
	return 0
}
