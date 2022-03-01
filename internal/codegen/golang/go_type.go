package golang

import (
	"github.com/kyleconroy/sqlc/internal/pattern"
	"github.com/kyleconroy/sqlc/internal/plugin"
)

// XXX: These are copied from python codegen.
func matchString(pat, target string) bool {
	matcher, err := pattern.MatchCompile(pat)
	if err != nil {
		panic(err)
	}
	return matcher.MatchString(target)
}

func matches(o *plugin.Override, n *plugin.Identifier, defaultSchema string) bool {
	if n == nil {
		return false
	}

	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}

	if o.Table.Catalog != "" && !matchString(o.Table.Catalog, n.Catalog) {
		return false
	}

	if o.Table.Schema == "" && schema != "" {
		return false
	}

	if o.Table.Schema != "" && !matchString(o.Table.Schema, schema) {
		return false
	}

	if o.Table.Name == "" && n.Name != "" {
		return false
	}

	if o.Table.Name != "" && !matchString(o.Table.Name, n.Name) {
		return false
	}

	return true
}

func goType(req *plugin.CodeGenRequest, col *plugin.Column) string {
	// Check if the column's type has been overridden
	for _, oride := range req.Settings.Overrides {
		if oride.GoType.TypeName == "" {
			continue
		}
		sameTable := matches(oride, col.Table, req.Catalog.DefaultSchema)
		if oride.Column != "" && matchString(oride.ColumnName, col.Name) && sameTable {
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
	columnType := dataType(col.Type)
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
