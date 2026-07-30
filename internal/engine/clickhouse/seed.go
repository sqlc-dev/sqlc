package clickhouse

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
)

func Dialect() core.Option {
	return core.WithSeed(Seed)
}

type chType struct {
	name     string
	category string
}

var clickhouseTypes = []chType{
	{"UInt8", "N"}, {"UInt16", "N"}, {"UInt32", "N"}, {"UInt64", "N"},
	{"UInt128", "N"}, {"UInt256", "N"},
	{"Int8", "N"}, {"Int16", "N"}, {"Int32", "N"}, {"Int64", "N"},
	{"Int128", "N"}, {"Int256", "N"},
	{"Float32", "N"}, {"Float64", "N"}, {"BFloat16", "N"},
	{"Decimal", "N"}, {"Decimal32", "N"}, {"Decimal64", "N"},
	{"Decimal128", "N"}, {"Decimal256", "N"},
	{"Bool", "B"},
	{"String", "S"}, {"FixedString", "S"}, {"UUID", "S"},
	{"Date", "D"}, {"Date32", "D"}, {"DateTime", "D"}, {"DateTime64", "D"},
	{"IPv4", "S"}, {"IPv6", "S"}, {"JSON", "U"},
	{"Enum8", "U"}, {"Enum16", "U"},
	{"Nullable", "U"}, {"LowCardinality", "U"}, {"Array", "A"},
	{"Map", "U"}, {"Tuple", "U"}, {"Nested", "U"}, {"Nothing", "U"},
}

func Seed(cat *core.Catalog) error {
	dialectOID, err := cat.CreateDialect("clickhouse")
	if err != nil {
		return err
	}
	for _, t := range clickhouseTypes {
		if _, err := cat.CreateTypeSpec(core.TypeSpec{
			Name:       t.name,
			Typtype:    "b",
			Category:   t.category,
			DialectOID: dialectOID,
		}); err != nil {
			return fmt.Errorf("seed clickhouse type %q: %w", t.name, err)
		}
	}
	return nil
}
