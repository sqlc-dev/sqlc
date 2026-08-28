package ast

type CommentStmt struct {
	Tag NodeTag[CommentStmt] `json:"tag"`

	Objtype ObjectType `json:"objtype"`
	Object  Node       `json:"object,omitempty"`
	Comment *string    `json:"comment,omitempty"`
}

func (n *CommentStmt) Pos() int {
	return 0
}
