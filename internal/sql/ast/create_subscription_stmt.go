package ast

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateSubscriptionStmt struct {
	Subname     *string
	Conninfo    *string
	Publication *List
	Options     *List
}

func (n *CreateSubscriptionStmt) Pos() int {
	return 0
}
