package dolphin

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// Dialect returns the catalog option that seeds MySQL's type system.
func Dialect() core.Option {
	return core.WithSeed(Seed)
}

var spec = seed.Spec{
	Dialect: "mysql",
	Types: []seed.Type{
		{Name: "bool", Category: "B", Aliases: []string{"boolean"}},

		{Name: "tinyint", Category: "N"},
		{Name: "smallint", Category: "N"},
		{Name: "mediumint", Category: "N"},
		{Name: "int", Category: "N", Aliases: []string{"integer"}},
		{Name: "bigint", Category: "N", Aliases: []string{"signed", "unsigned"}},
		{Name: "float", Category: "N"},
		{Name: "double", Category: "N", Aliases: []string{"double precision", "real"}},
		{Name: "decimal", Category: "N", Aliases: []string{"numeric", "dec", "fixed"}},
		{Name: "bit", Category: "N"},

		{Name: "char", Category: "S"},
		{Name: "varchar", Category: "S"},
		{Name: "tinytext", Category: "S"},
		{Name: "text", Category: "S"},
		{Name: "mediumtext", Category: "S"},
		{Name: "longtext", Category: "S"},
		{Name: "enum", Category: "S"},
		{Name: "set", Category: "S"},

		{Name: "binary", Category: "U"},
		{Name: "varbinary", Category: "U"},
		{Name: "tinyblob", Category: "U"},
		{Name: "blob", Category: "U"},
		{Name: "mediumblob", Category: "U"},
		{Name: "longblob", Category: "U"},
		{Name: "json", Category: "U"},
		{Name: "geometry", Category: "U"},

		{Name: "date", Category: "D"},
		{Name: "datetime", Category: "D"},
		{Name: "timestamp", Category: "D"},
		{Name: "time", Category: "D"},
		{Name: "year", Category: "D"},
	},
	Const: map[string]string{
		core.ConstInteger: "int",
		core.ConstFloat:   "decimal",
		core.ConstString:  "text",
		core.ConstBool:    "bool",
	},
	Bool:                 "bool",
	Comparison:           []string{"=", "<>", "!=", "<", "<=", ">", ">=", "<=>"},
	ComparisonCategories: "BNSDU",
	Arithmetic:           []string{"+", "-", "*", "/", "%", "DIV"},
	ArithmeticCategories: "N",
	CastCategories:       "NSD",
	Operators: []seed.Operator{
		{Name: "~~", Left: "text", Right: "text", Result: "bool"},
		{Name: "!~~", Left: "text", Right: "text", Result: "bool"},
		{Name: "REGEXP", Left: "text", Right: "text", Result: "bool"},
	},
}

func Seed(cat *core.Catalog) error {
	s := spec
	s.Funcs = defaultSchema("public").Funcs
	return s.Apply(cat)
}
