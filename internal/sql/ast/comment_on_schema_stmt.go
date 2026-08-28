package ast

type CommentOnSchemaStmt struct {
	Tag NodeTag[CommentOnSchemaStmt] `json:"tag"`

	Schema  *String `json:",omitempty"`
	Comment *string `json:",omitempty"`
}

func (n *CommentOnSchemaStmt) Pos() int {
	return 0
}
