package golang

import (
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func goType(req *plugin.CodeGenRequest, col *compiler.Column) string {
	// Check if the column's type has been overridden
	for _, oride := range req.Settings.Overrides {
		if oride.GoTypeName == "" {
			continue
		}
		sameTable := oride.Matches(col.Table, req.Catalog.DefaultSchema)
		if oride.Column != "" && oride.ColumnName.MatchString(col.Name) && sameTable {
			return oride.GoTypeName
		}
	}
	typ := goInnerType(req, col)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goInnerType(req *plugin.CodeGenRequest, col *compiler.Column) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, oride := range req.Settings.Overrides {
		if oride.GoTypeName == "" {
			continue
		}
		if oride.DBType != "" && oride.DBType == columnType && oride.Nullable != notNull {
			return oride.GoTypeName
		}
	}

	// TODO: Extend the engine interface to handle types
	switch req.Settings.Engine {
	case "mysql":
		return mysqlType(req, col)
	case "postgres":
		return postgresType(req, col)
	case "_lemon":
		return sqliteType(req, col)
	default:
		return "interface{}"
	}
}
