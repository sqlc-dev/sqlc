package sqlite

import (
	"github.com/kyleconroy/sqlc/internal/engine/sqlite/parser"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type tableNamer interface {
	Table_name() parser.ITable_nameContext
	Database_name() parser.IDatabase_nameContext
}

func parseTableName(c tableNamer) *ast.TableName {
	name := ast.TableName{
		Name: c.Table_name().GetText(),
	}
	if c.Database_name() != nil {
		name.Schema = c.Database_name().GetText()
	}
	return &name
}

func hasNotNullConstraint(checks []parser.IColumn_constraintContext) bool {
	for i := range checks {
		constraint, ok := checks[i].(*parser.Column_constraintContext)
		if !ok {
			continue
		}
		if constraint.K_PRIMARY() != nil && constraint.K_KEY() != nil {
			return true
		}
		if constraint.K_NOT() != nil && constraint.K_NULL() != nil {
			return true
		}
	}
	return false
}
