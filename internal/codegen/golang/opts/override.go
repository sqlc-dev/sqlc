package opts

import (
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type Override struct {
	CodeType   string             `json:"code_type"`
	DbType     string             `json:"db_type"`
	Nullable   bool               `json:"nullable"`
	Column     string             `json:"column"`
	Table      *plugin.Identifier `json:"table"`
	ColumnName string             `json:"column_name"`
	GoType     *ParsedGoType      `json:"go_type"`
	Unsigned   bool               `json:"unsigned"`
}

type ParsedGoType struct {
	ImportPath string            `json:"import_path"`
	Package    string            `json:"package"`
	TypeName   string            `json:"type_name"`
	BasicType  bool              `json:"basic_type"`
	StructTags map[string]string `json:"struct_tags"`
}

func (o *Override) Matches(n *plugin.Identifier, defaultSchema string) bool {
	if n == nil {
		return false
	}
	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}
	if o.Table.Catalog != "" && !sdk.MatchString(o.Table.Catalog, n.Catalog) {
		return false
	}
	if o.Table.Schema == "" && schema != "" {
		return false
	}
	if o.Table.Schema != "" && !sdk.MatchString(o.Table.Schema, schema) {
		return false
	}
	if o.Table.Name == "" && n.Name != "" {
		return false
	}
	if o.Table.Name != "" && !sdk.MatchString(o.Table.Name, n.Name) {
		return false
	}
	return true
}
