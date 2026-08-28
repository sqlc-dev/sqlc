package ast

type CommentOnColumnStmt struct {
	Tag NodeTag[CommentOnColumnStmt] `json:"tag"`

	Table   *TableName `json:",omitempty"`
	Col     *ColumnRef `json:",omitempty"`
	Comment *string    `json:",omitempty"`
}

func (n *CommentOnColumnStmt) Pos() int {
	return 0
}
