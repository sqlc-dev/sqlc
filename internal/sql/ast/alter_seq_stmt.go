package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type AlterSeqStmt struct {
	Sequence    *RangeVar
	Options     *ast.List
	ForIdentity bool
	MissingOk   bool
}

func (n *AlterSeqStmt) Pos() int {
	return 0
}
