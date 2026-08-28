package ast

type CommentOnTypeStmt struct {
	Tag NodeTag[CommentOnTypeStmt] `json:"tag"`

	Type    *TypeName `json:"type,omitempty"`
	Comment *string   `json:"comment,omitempty"`
}

func (n *CommentOnTypeStmt) Pos() int {
	return 0
}
