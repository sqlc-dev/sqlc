package postgresql

import (
	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// Dialect returns the catalog option that seeds PostgreSQL's type system.
func Dialect() core.Option {
	return core.WithSeed(Seed)
}

var spec = seed.Spec{
	Dialect: "postgresql",
	Types: []seed.Type{
		{Name: "bool", Category: "B", Aliases: []string{"boolean"}},

		{Name: "int2", Category: "N", Aliases: []string{"smallint"}},
		{Name: "int4", Category: "N", Aliases: []string{"integer", "int"}},
		{Name: "int8", Category: "N", Aliases: []string{"bigint"}},
		{Name: "float4", Category: "N", Aliases: []string{"real"}},
		{Name: "float8", Category: "N", Aliases: []string{"double precision"}},
		{Name: "numeric", Category: "N", Aliases: []string{"decimal"}},
		{Name: "money", Category: "N"},
		{Name: "oid", Category: "N"},
		// The serial types are column definition sugar for an integer plus a
		// sequence, but they are what a schema declares, so the catalog carries
		// them as types of their own.
		{Name: "serial2", Category: "N", Aliases: []string{"smallserial"}},
		{Name: "serial4", Category: "N", Aliases: []string{"serial"}},
		{Name: "serial8", Category: "N", Aliases: []string{"bigserial"}},

		{Name: "text", Category: "S"},
		{Name: "varchar", Category: "S", Aliases: []string{"character varying"}},
		{Name: "bpchar", Category: "S", Aliases: []string{"char", "character"}},
		{Name: "name", Category: "S"},
		{Name: "citext", Category: "S"},

		{Name: "date", Category: "D"},
		{Name: "time", Category: "D", Aliases: []string{"time without time zone"}},
		{Name: "timetz", Category: "D", Aliases: []string{"time with time zone"}},
		{Name: "timestamp", Category: "D", Aliases: []string{"timestamp without time zone"}},
		{Name: "timestamptz", Category: "D", Aliases: []string{"timestamp with time zone"}},

		{Name: "interval", Category: "T"},

		{Name: "uuid", Category: "U"},
		{Name: "bytea", Category: "U"},
		{Name: "json", Category: "U"},
		{Name: "jsonb", Category: "U"},
		{Name: "xml", Category: "U"},
		{Name: "inet", Category: "U"},
		{Name: "cidr", Category: "U"},
		{Name: "macaddr", Category: "U"},
		{Name: "macaddr8", Category: "U"},
		{Name: "bit", Category: "U"},
		{Name: "varbit", Category: "U", Aliases: []string{"bit varying"}},
		{Name: "tsvector", Category: "U"},
		{Name: "tsquery", Category: "U"},
		{Name: "point", Category: "U"},
		{Name: "line", Category: "U"},
		{Name: "lseg", Category: "U"},
		{Name: "box", Category: "U"},
		{Name: "path", Category: "U"},
		{Name: "polygon", Category: "U"},
		{Name: "circle", Category: "U"},
		{Name: "int4range", Category: "U"},
		{Name: "int8range", Category: "U"},
		{Name: "numrange", Category: "U"},
		{Name: "tsrange", Category: "U"},
		{Name: "tstzrange", Category: "U"},
		{Name: "daterange", Category: "U"},

		{Name: "anyarray", Category: "A", Aliases: []string{"array"}},
	},
	Const: map[string]string{
		core.ConstInteger: "int4",
		core.ConstFloat:   "numeric",
		core.ConstString:  "text",
		core.ConstBool:    "bool",
	},
	Bool:                 "bool",
	Comparison:           []string{"=", "<>", "!=", "<", "<=", ">", ">=", "IS DISTINCT FROM", "IS NOT DISTINCT FROM"},
	ComparisonCategories: "BNSDTU",
	Arithmetic:           []string{"+", "-", "*", "/", "%"},
	ArithmeticCategories: "N",
	// Every numeric spelling compares with every other, and so does every
	// string spelling, which is what an unadorned "WHERE bigint_col = 1"
	// depends on.
	CastCategories: "NSD",
	Operators: []seed.Operator{
		{Name: "||", Left: "text", Right: "text", Result: "text"},
		{Name: "~~", Left: "text", Right: "text", Result: "bool"},
		{Name: "!~~", Left: "text", Right: "text", Result: "bool"},
		{Name: "~~*", Left: "text", Right: "text", Result: "bool"},
		{Name: "!~~*", Left: "text", Right: "text", Result: "bool"},
		{Name: "~", Left: "text", Right: "text", Result: "bool"},
		{Name: "!~", Left: "text", Right: "text", Result: "bool"},
		{Name: "~*", Left: "text", Right: "text", Result: "bool"},
		{Name: "!~*", Left: "text", Right: "text", Result: "bool"},
		{Name: "@>", Left: "jsonb", Right: "jsonb", Result: "bool"},
		{Name: "<@", Left: "jsonb", Right: "jsonb", Result: "bool"},
		{Name: "->", Left: "jsonb", Right: "text", Result: "jsonb"},
		{Name: "->>", Left: "jsonb", Right: "text", Result: "text"},
		{Name: "->", Left: "json", Right: "text", Result: "json"},
		{Name: "->>", Left: "json", Right: "text", Result: "text"},
		{Name: "+", Left: "date", Right: "int4", Result: "date"},
		{Name: "-", Left: "date", Right: "int4", Result: "date"},
		{Name: "+", Left: "timestamptz", Right: "interval", Result: "timestamptz"},
		{Name: "-", Left: "timestamptz", Right: "interval", Result: "timestamptz"},
		{Name: "+", Left: "timestamp", Right: "interval", Result: "timestamp"},
		{Name: "-", Left: "timestamp", Right: "interval", Result: "timestamp"},
	},
}

func Seed(cat *core.Catalog) error {
	s := spec
	s.Funcs = genPGCatalog().Funcs
	return s.Apply(cat)
}
