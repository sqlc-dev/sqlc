package compiler

import (
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
)

type Result struct {
	Catalog *catalog.Catalog
	Queries []*Query
}
