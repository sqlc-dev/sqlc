package ast

import ()

type CommentStmt struct {
	Objtype ObjectType
	Object  Node
	Comment *string
}

func (n *CommentStmt) Pos() int {
	return 0
}
