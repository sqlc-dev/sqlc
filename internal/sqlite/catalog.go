package sqlite

import "github.com/kyleconroy/sqlc/internal/sql/catalog"

func NewCatalog() *catalog.Catalog {
	c := catalog.New("main")
	return c
}
