package oracle

import "github.com/sqlc-dev/sqlc/internal/sql/catalog"

// defaultSchemaName is the catalog's default schema. Oracle resolves unqualified
// objects against the connected user's schema; sqlc needs a concrete default, so
// we use a neutral name that CREATE TABLE statements without a schema prefix
// resolve into.
const defaultSchemaName = "public"

// NewCatalog returns a catalog seeded with Oracle's built-in schema, types and
// functions. Schema objects (tables, views) are added by the compiler as it
// walks the parsed DDL.
func NewCatalog() *catalog.Catalog {
	return &catalog.Catalog{
		DefaultSchema: defaultSchemaName,
		Schemas: []*catalog.Schema{
			defaultSchema(defaultSchemaName),
		},
		Extensions: map[string]struct{}{},
	}
}

// newTestCatalog returns an empty catalog for use in unit tests.
func newTestCatalog() *catalog.Catalog {
	return catalog.New(defaultSchemaName)
}
