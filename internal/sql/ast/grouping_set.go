package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type GroupingSet struct {
	Kind     GroupingSetKind
	Content  *List
	Location int
}

func (n *GroupingSet) Pos() int {
	return n.Location
}
