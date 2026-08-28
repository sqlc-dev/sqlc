package ast

type CommentOnTypeStmt struct {
	Tag NodeTag[CommentOnTypeStmt] `json:"tag"`

	Type    *TypeName `json:",omitempty"`
	Comment *string   `json:",omitempty"`
}

func (n *CommentOnTypeStmt) Pos() int {
	return 0
}
