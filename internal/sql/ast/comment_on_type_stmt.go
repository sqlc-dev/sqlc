package ast

type CommentOnTypeStmt struct {
	Tag NodeTag[CommentOnTypeStmt] `json:"tag"`

	Type    *TypeName
	Comment *string
}

func (n *CommentOnTypeStmt) Pos() int {
	return 0
}
