package ast

type CommentOnSchemaStmt struct {
	Tag NodeTag[CommentOnSchemaStmt] `json:"tag"`

	Schema  *String `json:"schema,omitempty"`
	Comment *string `json:"comment,omitempty"`
}

func (n *CommentOnSchemaStmt) Pos() int {
	return 0
}
