package postgresql

import "github.com/kyleconroy/sqlc/internal/sql/catalog"

func NewCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, pgTemp())
	c.Schemas = append(c.Schemas, pgCatalog())
	c.SearchPath = []string{"pg_catalog"}
	return c
}
