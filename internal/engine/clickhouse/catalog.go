package clickhouse

import "github.com/sqlc-dev/sqlc/internal/sql/catalog"

// NewCatalog creates a new catalog for ClickHouse.
// ClickHouse uses "default" as the default database/schema.
func NewCatalog() *catalog.Catalog {
	def := "default"
	return &catalog.Catalog{
		DefaultSchema: def,
		Schemas: []*catalog.Schema{
			defaultSchema(def),
		},
		Extensions: map[string]struct{}{},
	}
}

// defaultSchema creates a default schema with ClickHouse system tables.
func defaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{
		Name: name,
	}
}
