package sqlite

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// Dialect returns the catalog option that seeds SQLite's type system.
func Dialect() core.Option {
	return core.WithSeed(Seed)
}

// SQLite has five storage classes and assigns a column one of them from the
// declared type's affinity, but a schema may declare any name at all. The
// types below are the spellings the affinity rules call out, grouped by the
// storage class they land in.
var spec = seed.Spec{
	Dialect: "sqlite",
	Types: []seed.Type{
		{Name: "integer", Category: "N", Aliases: []string{
			"int", "tinyint", "smallint", "mediumint", "bigint",
			"unsigned big int", "int2", "int8",
		}},
		{Name: "real", Category: "N", Aliases: []string{"double", "double precision", "float"}},
		{Name: "numeric", Category: "N", Aliases: []string{"decimal"}},
		{Name: "boolean", Category: "B", Aliases: []string{"bool"}},
		{Name: "text", Category: "S", Aliases: []string{
			"character", "varchar", "varying character", "nchar",
			"native character", "nvarchar", "clob",
		}},
		{Name: "blob", Category: "U"},
		{Name: "any", Category: "U"},
		{Name: "date", Category: "D"},
		{Name: "datetime", Category: "D"},
		{Name: "timestamp", Category: "D"},
	},
	Const: map[string]string{
		core.ConstInteger: "integer",
		core.ConstFloat:   "real",
		core.ConstString:  "text",
		core.ConstBool:    "boolean",
	},
	Bool:                 "boolean",
	Comparison:           []string{"=", "==", "<>", "!=", "<", "<=", ">", ">=", "IS", "IS NOT"},
	ComparisonCategories: "BNSDU",
	Arithmetic:           []string{"+", "-", "*", "/", "%"},
	ArithmeticCategories: "N",
	// SQLite compares values of any declared type, so every type it knows is
	// implicitly castable to every other.
	CastCategories: "*",
	Operators: []seed.Operator{
		{Name: "||", Left: "text", Right: "text", Result: "text"},
		{Name: "~~", Left: "text", Right: "text", Result: "boolean"},
		{Name: "!~~", Left: "text", Right: "text", Result: "boolean"},
	},
}

func Seed(cat *core.Catalog) error {
	s := spec
	s.Funcs = defaultSchema("main").Funcs
	return s.Apply(cat)
}
