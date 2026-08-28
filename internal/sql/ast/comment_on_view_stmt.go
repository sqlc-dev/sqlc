package ast

type CommentOnViewStmt struct {
	Tag NodeTag[CommentOnViewStmt] `json:"tag"`

	View    *TableName
	Comment *string
}

func (n *CommentOnViewStmt) Pos() int {
	return 0
}
