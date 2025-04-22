package ydb

import (
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{Name: name}
	s.Funcs = []*catalog.Function{}

	return s
}
