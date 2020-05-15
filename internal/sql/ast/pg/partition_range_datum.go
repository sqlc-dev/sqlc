package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PartitionRangeDatum struct {
	Kind     PartitionRangeDatumKind
	Value    ast.Node
	Location int
}

func (n *PartitionRangeDatum) Pos() int {
	return n.Location
}
