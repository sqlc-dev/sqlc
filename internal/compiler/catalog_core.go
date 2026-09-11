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
				// Codegen reads a data type and an array flag, and renders
				// a "[]" per dimension, so an array of arrays of integers
				// is the type integer with two dimensions.
				expr, err := c.TypeExprOf(col.TypeOID)
				if err != nil {
					return nil, err
				}
				inner := expr.Innermost()
				column := &catalog.Column{
					Name:       col.Name,
					Type:       ast.TypeName{Name: inner.Name},
					IsNotNull:  col.NotNull,
					IsArray:    expr.IsArray(),
					ArrayDims:  expr.ArrayDims(),
					IsUnsigned: strings.HasSuffix(inner.Name, " unsigned"),
				}
				if len(inner.Args) > 0 && inner.Args[0].Int != nil {
					l := int(*inner.Args[0].Int)
					column.Length = &l
				}
				t.Columns = append(t.Columns, column)
			}
			schema.Tables = append(schema.Tables, t)
		}
		cat.Schemas = append(cat.Schemas, schema)
	}
	return cat, nil
}
