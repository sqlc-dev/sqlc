package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterTSConfigurationStmt struct {
	Kind      AlterTSConfigType
	Cfgname   *ast.List
	Tokentype *ast.List
	Dicts     *ast.List
	Override  bool
	Replace   bool
	MissingOk bool
}

func (n *AlterTSConfigurationStmt) Pos() int {
	return 0
}
