package golang

import (
	"github.com/kyleconroy/sqlc/internal/codegen/sdk"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

func goType(req *plugin.CodeGenRequest, col *plugin.Column) string {
	// Check if the column's type has been overridden
	for _, oride := range req.Settings.Overrides {
		if oride.GoType.TypeName == "" {
			continue
		}
		sameTable := sdk.Matches(oride, col.Table, req.Catalog.DefaultSchema)
		if oride.Column != "" && sdk.MatchString(oride.ColumnName, col.Name) && sameTable {
			return oride.GoType.TypeName
		}
	}
	typ := goInnerType(req, col)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goInnerType(req *plugin.CodeGenRequest, col *plugin.Column) string {
	columnType := sdk.DataType(col.Type)
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, oride := range req.Settings.Overrides {
		if oride.GoType.TypeName == "" {
			continue
		}
		if oride.DbType != "" && oride.DbType == columnType && oride.Nullable != notNull {
			return oride.GoType.TypeName
		}
	}

	// TODO: Extend the engine interface to handle types
	switch req.Settings.Engine {
	case "mysql":
		return mysqlType(req, col)
	case "postgresql":
		return postgresType(req, col)
	case "_lemon":
		return sqliteType(req, col)
	default:
		return "interface{}"
	}
}
