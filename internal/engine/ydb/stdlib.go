package ydb

import (
	"github.com/sqlc-dev/sqlc/internal/engine/ydb/lib"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

func defaultSchema(name string) *catalog.Schema {
	s := &catalog.Schema{
		Name:  name,
		Funcs: []*catalog.Function{},
	}

	s.Funcs = append(s.Funcs, lib.BasicFunctions()...)
	s.Funcs = append(s.Funcs, lib.AggregateFunctions()...)
	s.Funcs = append(s.Funcs, lib.WindowFunctions()...)
	s.Funcs = append(s.Funcs, lib.CppFunctions()...)
	// TODO: add container functions if

	return s
}
