package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PartitionBoundSpec struct {
	Strategy    byte
	Listdatums  *ast.List
	Lowerdatums *ast.List
	Upperdatums *ast.List
	Location    int
}

func (n *PartitionBoundSpec) Pos() int {
	return n.Location
}
