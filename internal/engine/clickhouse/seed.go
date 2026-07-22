package clickhouse

import (
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/core"
)

// Dialect returns a core.Option that registers the "clickhouse" dialect
// and seeds the catalog with the built-in ClickHouse data types. It is
// the ClickHouse counterpart to the per-dialect Dialect() helpers in
// xqlc, and is how a ClickHouse-backed core.Catalog is constructed:
//
//	cat, err := core.New(clickhouse.Dialect())
func Dialect() core.Option {
	return core.WithSeed(Seed)
}

// chType pairs a ClickHouse type name with its pg_type-style category so
// the catalog can reason about it during (future) operator and cast
// resolution. Names are stored lower-cased by the catalog.
type chType struct {
	name     string
	category string // 'N'umeric | 'S'tring | 'B'oolean | 'D'atetime | 'A'rray | 'U'serdef
}

// clickhouseTypes is the built-in type surface. It is intentionally a
// curated hand list for now; once testgen/clickhouse exists these will
// be regenerated against a real ClickHouse instance, matching the
// seedgen approach used by the other engines.
var clickhouseTypes = []chType{
	// Unsigned integers
	{"UInt8", "N"}, {"UInt16", "N"}, {"UInt32", "N"}, {"UInt64", "N"},
	{"UInt128", "N"}, {"UInt256", "N"},
	// Signed integers
	{"Int8", "N"}, {"Int16", "N"}, {"Int32", "N"}, {"Int64", "N"},
	{"Int128", "N"}, {"Int256", "N"},
	// Floating point / decimal
	{"Float32", "N"}, {"Float64", "N"}, {"BFloat16", "N"},
	{"Decimal", "N"}, {"Decimal32", "N"}, {"Decimal64", "N"},
	{"Decimal128", "N"}, {"Decimal256", "N"},
	// Boolean
	{"Bool", "B"},
	// Strings / binary
	{"String", "S"}, {"FixedString", "S"}, {"UUID", "S"},
	// Date & time
	{"Date", "D"}, {"Date32", "D"}, {"DateTime", "D"}, {"DateTime64", "D"},
	// Network / misc scalar
	{"IPv4", "S"}, {"IPv6", "S"}, {"JSON", "U"},
	// Enums are declared with their variants; the base names still resolve.
	{"Enum8", "U"}, {"Enum16", "U"},
	// Wrapper / composite type names. Columns declared with these are
	// normally unwrapped to their inner scalar by the DDL handler, but we
	// register the names so an un-unwrappable declaration still resolves.
	{"Nullable", "U"}, {"LowCardinality", "U"}, {"Array", "A"},
	{"Map", "U"}, {"Tuple", "U"}, {"Nested", "U"}, {"Nothing", "U"},
}

// Seed registers the ClickHouse dialect row and its built-in types on a
// freshly bootstrapped catalog.
func Seed(cat *core.Catalog) error {
	dialectOID, err := cat.CreateDialect("clickhouse")
	if err != nil {
		return err
	}
	// ClickHouse identifiers are case-sensitive, so unlike MySQL/SQLite we
	// do not set a case-folding flag.
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
