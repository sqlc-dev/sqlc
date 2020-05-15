package pg

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type RuleStmt struct {
	Relation    *RangeVar
	Rulename    *string
	WhereClause ast.Node
	Event       CmdType
	Instead     bool
	Actions     *ast.List
	Replace     bool
}

func (n *RuleStmt) Pos() int {
	return 0
}
