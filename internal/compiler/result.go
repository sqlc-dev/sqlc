package compiler

import (
	"github.com/Dionid/sqlc/internal/sql/catalog"
)

type Result struct {
	Catalog *catalog.Catalog
	Queries []*Query
}
