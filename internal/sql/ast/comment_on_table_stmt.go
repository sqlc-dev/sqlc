package ast

type CommentOnTableStmt struct {
	Table   *TableName
	Comment *string
}

func (n *CommentOnTableStmt) Pos() int {
	return 0
}
