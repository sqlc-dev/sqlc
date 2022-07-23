package zetasql

import (
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

func NewCatalog() *catalog.Catalog {
	def := "public" // TODO: What is the default database for ZetaSQL?
	return &catalog.Catalog{
		DefaultSchema: def,
		Schemas:       []*catalog.Schema{},
		Extensions:    map[string]struct{}{},
	}
}
