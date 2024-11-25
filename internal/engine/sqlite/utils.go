package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite/parser"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type tableNamer interface {
	Table_name() parser.ITable_nameContext
	Schema_name() parser.ISchema_nameContext
}

func parseTableName(c tableNamer) *ast.TableName {
	name := ast.TableName{
		Name: identifier(c.Table_name().GetText()),
	}
	if c.Schema_name() != nil {
		name.Schema = c.Schema_name().GetText()
	}
	return &name
}

func hasNotNullConstraint(checks []parser.IColumn_constraintContext) bool {
	for i := range checks {
		constraint, ok := checks[i].(*parser.Column_constraintContext)
		if !ok {
			continue
		}
		if constraint.PRIMARY_() != nil && constraint.KEY_() != nil {
			return true
		}
		if constraint.NOT_() != nil && constraint.NULL_() != nil {
			return true
		}
	}
	return false
}
