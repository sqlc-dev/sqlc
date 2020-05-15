package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type DeclareCursorStmt struct {
	Portalname *string
	Options    int
	Query      ast.Node
}

func (n *DeclareCursorStmt) Pos() int {
	return 0
}
