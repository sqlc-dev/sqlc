package dolphin

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func NewCatalog() *catalog.Catalog {
	def := "public" // TODO: What is the default database for MySQL?
	return &catalog.Catalog{
		DefaultSchema: def,
		Schemas: []*catalog.Schema{
			defaultSchema(def),
		},
		Extensions: map[string]struct{}{},
	}
}
