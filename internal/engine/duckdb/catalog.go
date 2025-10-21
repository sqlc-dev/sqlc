package duckdb

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// NewCatalog creates a minimal catalog for DuckDB
// DuckDB uses database-backed catalog via analyzer
// This catalog is minimal - all type information comes from the database
func NewCatalog() *catalog.Catalog {
	def := "main"
	return &catalog.Catalog{
		DefaultSchema: def,
		Name:          "memory", // DuckDB's default catalog
		Schemas:       []*catalog.Schema{},
		SearchPath:    []string{def},
		Extensions:    map[string]struct{}{},
	}
}
