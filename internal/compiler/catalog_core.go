package compiler

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// coreResultCatalog dumps the core catalog into the legacy catalog shape a
// Result carries, so codegen sees the same table models either way a query
// set was analyzed. Only relations make the trip: codegen reads tables and
// their columns to build models, and none of the types, functions or
// operators the core catalog also holds.
func coreResultCatalog(c *core.Catalog) (*catalog.Catalog, error) {
	cat := catalog.New("public")
	namespaces, err := c.Namespaces()
	if err != nil {
		return nil, err
	}
	for _, ns := range namespaces {
		schema := &catalog.Schema{Name: ns.Name}
		tables, err := c.TablesInNamespace(ns.OID)
		if err != nil {
			return nil, err
		}
		for _, table := range tables {
			cols, err := c.ClassCodegenColumns(table.OID)
			if err != nil {
				return nil, err
			}
			t := &catalog.Table{Rel: &ast.TableName{Schema: ns.Name, Name: table.Name}}
			for _, col := range cols {
				// The catalog names an array type after its element with the
				// suffix appended, which is codegen's data type and array
				// flag in one string. The core catalog holds one dimension,
				// and codegen renders a "[]" per dimension.
				dataType, isArray := strings.CutSuffix(col.TypeName, core.ArraySuffix)
				column := &catalog.Column{
					Name:      col.Name,
					Type:      ast.TypeName{Name: dataType},
					IsNotNull: col.NotNull,
					IsArray:   isArray,
				}
				if isArray {
					column.ArrayDims = 1
				}
				t.Columns = append(t.Columns, column)
			}
			schema.Tables = append(schema.Tables, t)
		}
		cat.Schemas = append(cat.Schemas, schema)
	}
	return cat, nil
}
