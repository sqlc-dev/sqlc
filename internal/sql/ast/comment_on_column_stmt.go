package ast

type CommentOnColumnStmt struct {
	Tag NodeTag[CommentOnColumnStmt] `json:"tag"`

	Table   *TableName `json:"table,omitempty"`
	Col     *ColumnRef `json:"col,omitempty"`
	Comment *string    `json:"comment,omitempty"`
}

func (n *CommentOnColumnStmt) Pos() int {
	return 0
}
