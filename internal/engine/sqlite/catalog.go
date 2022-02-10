package sqlite

import "github.com/egtann/sqlc/internal/sql/catalog"

func NewCatalog() *catalog.Catalog {
	c := catalog.New("main")
	return c
}
