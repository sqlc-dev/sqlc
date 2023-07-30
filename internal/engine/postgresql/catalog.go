package postgresql

import "github.com/sqlc-dev/sqlc/internal/sql/catalog"

// toPointer converts an int to a pointer without a temporary
// variable at the call-site, and is used by the generated schemas
func toPointer(x int) *int {
	return &x
}

func NewCatalog() *catalog.Catalog {
	c := catalog.New("public")
	c.Schemas = append(c.Schemas, pgTemp())
	c.Schemas = append(c.Schemas, genPGCatalog())
	c.Schemas = append(c.Schemas, genInformationSchema())
	c.SearchPath = []string{"pg_catalog"}
	c.LoadExtension = loadExtension
	return c
}
