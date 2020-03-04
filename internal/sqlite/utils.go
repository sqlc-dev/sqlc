package sqlite

import (
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sqlite/parser"
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
