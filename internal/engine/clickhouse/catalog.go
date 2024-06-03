package clickhouse

import "github.com/sqlc-dev/sqlc/internal/sql/catalog"

func NewCatalog() *catalog.Catalog {
	def := ""
	c := &catalog.Catalog{
		DefaultSchema: def,
		Schemas: []*catalog.Schema{
			&catalog.Schema{Name: ""},
		},
		Extensions: map[string]struct{}{},
	}
	c.Schemas = append(c.Schemas, getInformationSchema())
	return c
}
