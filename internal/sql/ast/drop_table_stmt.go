package ast

type DropTableStmt struct {
	Tag NodeTag[DropTableStmt] `json:"tag"`

	IfExists bool         `json:"if_exists"`
	Tables   []*TableName `json:"tables,omitempty"`
}

func (n *DropTableStmt) Pos() int {
	return 0
}
