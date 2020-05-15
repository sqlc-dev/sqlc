package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type CreateTableAsStmt struct {
	Query        ast.Node
	Into         *IntoClause
	Relkind      ObjectType
	IsSelectInto bool
	IfNotExists  bool
}

func (n *CreateTableAsStmt) Pos() int {
	return 0
}
