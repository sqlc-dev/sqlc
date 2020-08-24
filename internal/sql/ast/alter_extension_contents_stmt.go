package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterExtensionContentsStmt struct {
	Extname *string
	Action  int
	Objtype ObjectType
	Object  ast.Node
}

func (n *AlterExtensionContentsStmt) Pos() int {
	return 0
}
