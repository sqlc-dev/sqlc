package ast

type CommentOnTableStmt struct {
	Tag NodeTag[CommentOnTableStmt] `json:"tag"`

	Table   *TableName `json:",omitempty"`
	Comment *string    `json:",omitempty"`
}

func (n *CommentOnTableStmt) Pos() int {
	return 0
}
