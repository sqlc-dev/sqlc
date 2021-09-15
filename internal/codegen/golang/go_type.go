package golang

import (
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
)

func goType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	// Check if the column's type has been overridden
	for _, oride := range settings.Overrides {
		if oride.GoTypeName == "" {
			continue
		}
		sameTable := oride.Matches(col.Table, r.Catalog.DefaultSchema)
		if oride.Column != "" && oride.ColumnName.MatchString(col.Name) && sameTable {
			return oride.GoTypeName
		}
	}
	typ := goInnerType(r, col, settings)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func goInnerType(r *compiler.Result, col *compiler.Column, settings config.CombinedSettings) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, oride := range settings.Overrides {
		if oride.GoTypeName == "" {
			continue
		}
		if oride.DBType != "" && oride.DBType == columnType && oride.Nullable != notNull {
			return oride.GoTypeName
		}
	}

	// TODO: Extend the engine interface to handle types
	switch settings.Package.Engine {
	case config.EngineMySQL:
		return mysqlType(r, col, settings)
	case config.EnginePostgreSQL:
		return postgresType(r, col, settings)
	case config.EngineXLemon:
		return sqliteType(r, col, settings)
	default:
		return "interface{}"
	}
}
