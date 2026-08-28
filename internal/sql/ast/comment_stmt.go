package ast

type CommentStmt struct {
	Tag NodeTag[CommentStmt] `json:"tag"`

	Objtype ObjectType
	Object  Node    `json:",omitempty"`
	Comment *string `json:",omitempty"`
}

func (n *CommentStmt) Pos() int {
	return 0
}
