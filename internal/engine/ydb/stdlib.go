package ydb

import (
	"github.com/sqlc-dev/sqlc/internal/engine/ydb/lib"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{
		Name:  name,
		Funcs: make([]*catalog.Function, 0, 128),
	}

	s.Funcs = append(s.Funcs, lib.BasicFunctions()...)
	s.Funcs = append(s.Funcs, lib.AggregateFunctions()...)

	return s
}
