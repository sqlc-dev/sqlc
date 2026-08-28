package ast

type CommentOnSchemaStmt struct {
	Tag NodeTag[CommentOnSchemaStmt] `json:"tag"`

	Schema  *String
	Comment *string
}

func (n *CommentOnSchemaStmt) Pos() int {
	return 0
}
