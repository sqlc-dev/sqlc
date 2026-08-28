package ast

type CommentOnColumnStmt struct {
	Tag NodeTag[CommentOnColumnStmt] `json:"tag"`

	Table   *TableName
	Col     *ColumnRef
	Comment *string
}

func (n *CommentOnColumnStmt) Pos() int {
	return 0
}
