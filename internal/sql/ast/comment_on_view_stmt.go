package ast

type CommentOnViewStmt struct {
	Tag NodeTag[CommentOnViewStmt] `json:"tag"`

	View    *TableName `json:",omitempty"`
	Comment *string    `json:",omitempty"`
}

func (n *CommentOnViewStmt) Pos() int {
	return 0
}
