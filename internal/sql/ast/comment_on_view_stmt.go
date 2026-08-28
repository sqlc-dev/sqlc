package ast

type CommentOnViewStmt struct {
	Tag NodeTag[CommentOnViewStmt] `json:"tag"`

	View    *TableName `json:"view,omitempty"`
	Comment *string    `json:"comment,omitempty"`
}

func (n *CommentOnViewStmt) Pos() int {
	return 0
}
