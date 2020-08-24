package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AccessPriv struct {
	PrivName *string
	Cols     *List
}

func (n *AccessPriv) Pos() int {
	return 0
}
