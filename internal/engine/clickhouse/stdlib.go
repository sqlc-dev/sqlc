package clickhouse

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	return &catalog.Schema{Name: name}
}
