package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CommentStmt struct {
	Objtype ObjectType
	Object  Node
	Comment *string
}

func (n *CommentStmt) Pos() int {
	return 0
}
