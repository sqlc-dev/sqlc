package ast

type CommentOnTypeStmt struct {
	Type    *TypeName
	Comment *string
}

func (n *CommentOnTypeStmt) Pos() int {
	return 0
}
