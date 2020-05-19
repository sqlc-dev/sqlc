package ast

type CommentOnSchemaStmt struct {
	Schema  *String
	Comment *string
}

func (n *CommentOnSchemaStmt) Pos() int {
	return 0
}
