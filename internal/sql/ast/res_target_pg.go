package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type ResTarget struct {
	Name        *string
	Indirection *List
	Val         Node
	Location    int
}

func (n *ResTarget) Pos() int {
	return n.Location
}
