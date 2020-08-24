package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type PartitionBoundSpec struct {
	Strategy    byte
	Listdatums  *List
	Lowerdatums *List
	Upperdatums *List
	Location    int
}

func (n *PartitionBoundSpec) Pos() int {
	return n.Location
}
