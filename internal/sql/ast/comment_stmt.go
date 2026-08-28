package ast

type CommentStmt struct {
	Tag NodeTag[CommentStmt] `json:"tag"`

	Objtype ObjectType
	Object  Node
	Comment *string
}

func (n *CommentStmt) Pos() int {
	return 0
}
