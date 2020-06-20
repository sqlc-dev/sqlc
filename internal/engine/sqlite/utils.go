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
