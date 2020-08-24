package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PartitionSpec struct {
	Strategy   *string
	PartParams *ast.List
	Location   int
}

func (n *PartitionSpec) Pos() int {
	return n.Location
}
