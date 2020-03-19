package postgresql

import "github.com/kyleconroy/sqlc/internal/sql/catalog"

func NewCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, &catalog.Schema{Name: "pg_temp"})
	return c
}
