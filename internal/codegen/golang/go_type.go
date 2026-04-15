package golang

import (
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
			for k, v := range oride.GoType.StructTags {
				tags[k] = v
			}
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
		for k, v := range oride.GoType.StructTags {
			tags[k] = v
		}
	}
}

func goType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
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
	typ := goInnerType(req, options, col)
	if col.IsSqlcSlice {
		return "[]" + typ
	}
	if col.IsArray {
		return strings.Repeat("[]", int(col.ArrayDims)) + typ
	}
	return typ
}

func goInnerType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	// package overrides have a higher precedence
	for _, override := range options.Overrides {
		oride := override.ShimOverride
		if oride.GoType.TypeName == "" {
			continue
		}
		if override.MatchesColumn(col) {
			return oride.GoType.TypeName
		}
	}

	// TODO: Extend the engine interface to handle types
	switch req.Settings.Engine {
	case "mysql":
		return mysqlType(req, options, col)
	case "postgresql":
		return postgresType(req, options, col)
	case "sqlite":
		return sqliteType(req, options, col)
	default:
		return "interface{}"
	}
}
