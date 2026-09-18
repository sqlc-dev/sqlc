package golang

import (
	"maps"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

func addExtraGoStructTags(tags map[string]string, req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) {
	for _, override := range options.Overrides {
		oride := override.ShimOverride
		if oride.GoType.StructTags == nil {
			continue
		}
		if override.MatchesColumn(col) {
			maps.Copy(tags, oride.GoType.StructTags)
			continue
		}
		if !override.Matches(col.Table, req.Catalog.DefaultSchema) {
			// Different table.
			continue
		}
		cname := col.Name
		if col.OriginalName != "" {
			cname = col.OriginalName
		}
		if !sdk.MatchString(oride.ColumnName, cname) {
			// Different column.
			continue
		}
		// Add the extra tags.
		maps.Copy(tags, oride.GoType.StructTags)
	}
}

// goType is the Go type of a column's value: what a query's row scans
// into, or what a table's model holds.
func goType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	return goValueType(req, options, col, false)
}

// goParamType is the Go type of a parameter's value: what is passed to the
// driver for it. It differs from goType only where a driver reads a value
// through a type it will not accept as an argument.
func goParamType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	return goValueType(req, options, col, true)
}

func goValueType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column, param bool) string {
	// Check if the column's type has been overridden
	for _, override := range options.Overrides {
		oride := override.ShimOverride

		if oride.GoType.TypeName == "" {
			continue
		}
		cname := col.Name
		if col.OriginalName != "" {
			cname = col.OriginalName
		}
		sameTable := override.Matches(col.Table, req.Catalog.DefaultSchema)
		if oride.Column != "" && sdk.MatchString(oride.ColumnName, cname) && sameTable {
			if col.IsSqlcSlice {
				return "[]" + oride.GoType.TypeName
			}
			return oride.GoType.TypeName
		}
	}
	typ, complete := goInnerType(req, options, col, param)
	if col.IsSqlcSlice {
		return "[]" + typ
	}
	if col.IsArray && !complete {
		return strings.Repeat("[]", int(col.ArrayDims)) + typ
	}
	return typ
}

// goInnerType is the Go type of the column's value. The legacy mappers
// return the element type of an array column and leave the dimensions to
// goType; the mappers that read the type expression render the whole type,
// arrays included, and say so with complete.
func goInnerType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column, param bool) (typ string, complete bool) {
	// package overrides have a higher precedence
	for _, override := range options.Overrides {
		oride := override.ShimOverride
		if oride.GoType.TypeName == "" {
			continue
		}
		if override.MatchesColumn(col) {
			return oride.GoType.TypeName, false
		}
	}

	// TODO: Extend the engine interface to handle types
	switch req.Settings.Engine {
	case engineMySQL:
		return mysqlType(req, options, col), false
	case enginePostgreSQL:
		return postgresType(req, options, col), false
	case engineSQLite:
		return sqliteType(req, options, col), false
	case engineClickHouse:
		return clickhouseType(req, options, col), true
	case engineDuckDB:
		return duckdbType(req, options, col, param), true
	case engineGoogleSQL:
		return googlesqlType(req, options, col), true
	case engineMSSQL:
		return mssqlType(req, options, col), true
	default:
		return "any", false
	}
}
