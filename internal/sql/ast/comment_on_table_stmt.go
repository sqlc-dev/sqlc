package ast

type CommentOnTableStmt struct {
	Tag NodeTag[CommentOnTableStmt] `json:"tag"`

	Table   *TableName
	Comment *string
}

func (n *CommentOnTableStmt) Pos() int {
	return 0
}
