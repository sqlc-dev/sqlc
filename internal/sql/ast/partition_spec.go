package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PartitionSpec struct {
	Strategy   *string
	PartParams *List
	Location   int
}

func (n *PartitionSpec) Pos() int {
	return n.Location
}
