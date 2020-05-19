package ast

type CommentOnColumnStmt struct {
	Table   *TableName
	Col     *ColumnRef
	Comment *string
}

func (n *CommentOnColumnStmt) Pos() int {
	return 0
}
