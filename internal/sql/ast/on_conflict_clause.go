package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type OnConflictClause struct {
	Action      OnConflictAction
	Infer       *InferClause
	TargetList  *List
	WhereClause Node
	Location    int
}

func (n *OnConflictClause) Pos() int {
	return n.Location
}
