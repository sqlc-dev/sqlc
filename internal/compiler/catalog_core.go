package compiler

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// coreResultCatalog dumps the core catalog into the legacy catalog shape a
// Result carries, so codegen sees the same table models either way a query
// set was analyzed. Relations and enums make the trip: codegen reads tables
// and views with their columns to build models and enums to build their Go
// types, and none of the functions or operators the core catalog also holds.
func coreResultCatalog(c *core.Catalog) (*catalog.Catalog, error) {
	cat := catalog.New("public")
	namespaces, err := c.Namespaces()
	if err != nil {
		return nil, err
	}
	schemas := map[string]*catalog.Schema{}
	for _, ns := range namespaces {
		schema := &catalog.Schema{Name: ns.Name}
		schemas[ns.Name] = schema
		tables, err := c.ModelClassesInNamespace(ns.OID)
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

	// The catalog spells an enum's schema in its name, and a schema that
	// holds only types has no namespace of its own, so one is made here.
	enums, err := c.Enums()
	if err != nil {
		return nil, err
	}
	for _, enum := range enums {
		labels, err := c.EnumLabels(enum.OID)
		if err != nil {
			return nil, err
		}
		schemaName, name := core.SplitTypeName(enum.Name)
		if schemaName == "" {
			schemaName = cat.DefaultSchema
		}
		schema, ok := schemas[schemaName]
		if !ok {
			schema = &catalog.Schema{Name: schemaName}
			schemas[schemaName] = schema
			cat.Schemas = append(cat.Schemas, schema)
		}
		schema.Types = append(schema.Types, &catalog.Enum{Name: name, Vals: labels})
	}
	return cat, nil
}
