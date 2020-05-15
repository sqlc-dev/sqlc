package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterSubscriptionStmt struct {
	Kind        AlterSubscriptionType
	Subname     *string
	Conninfo    *string
	Publication *ast.List
	Options     *ast.List
}

func (n *AlterSubscriptionStmt) Pos() int {
	return 0
}
