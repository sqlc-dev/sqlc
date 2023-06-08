package ast

type CommentOnViewStmt struct {
	View    *TableName
	Comment *string
}

func (n *CommentOnViewStmt) Pos() int {
	return 0
}
