package clickhouse

import (
	"embed"
	"sort"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/core/seed"
)

// The dialect directory describes ClickHouse's type system. types.jsonl is
// generated from system.data_type_families of the pinned ClickHouse release
// by goldeneye (internal/goldeneye), which also checks it against one;
// dialect.json and functions.jsonl are authored by hand, since ClickHouse
// publishes no function signatures. Regenerate from internal/goldeneye with:
//
//	go run ./cmd/goldeneye install clickhouse
//	go run ./cmd/goldeneye generate clickhouse
//
//go:embed dialect
var dialectFS embed.FS

// Dialect returns the catalog option that seeds ClickHouse's type system.
func Dialect() core.Option {
	return seed.Dialect(dialectFS, "dialect")
}

func init() {
	core.RegisterCanonicalizer("clickhouse", canonicalize)
	core.RegisterResultType("clickhouse", resultType)
}

// decimalPrecisions is the precision each sized decimal family stands for:
// ClickHouse stores Decimal32(s) as Decimal(9, s).
var decimalPrecisions = map[string]int64{
	"decimal32":  9,
	"decimal64":  18,
	"decimal128": 38,
	"decimal256": 76,
}

// canonicalize rewrites a type the way ClickHouse stores and reports it: a
// sized decimal is a Decimal with that precision, an Enum is an Enum8 or
// Enum16 with its members numbered, a Variant's members are sorted, and
// the function an aggregate-function type names is a word rather than a
// type.
func canonicalize(t *core.TypeExpr) *core.TypeExpr {
	switch t.Name {
	case "decimal32", "decimal64", "decimal128", "decimal256":
		if len(t.Args) == 1 && t.Args[0].Int != nil {
			p := decimalPrecisions[t.Name]
			return &core.TypeExpr{Name: "decimal", Nullable: t.Nullable, Args: []core.TypeArg{{Int: &p}, t.Args[0]}}
		}
	case "enum", "enum8", "enum16":
		out := t.Clone()
		if out.Name == "enum" {
			out.Name = "enum8"
			if len(out.Args) > 127 {
				out.Name = "enum16"
			}
		}
		next := int64(1)
		for i := range out.Args {
			a := &out.Args[i]
			switch {
			case a.Label != "" && a.Int != nil:
				next = *a.Int + 1
			case a.String != nil:
				// A bare member is numbered after the one before it.
				label := *a.String
				n := next
				*a = core.TypeArg{Label: label, Int: &n}
				next++
			}
		}
		return out
	case "variant":
		out := t.Clone()
		sort.SliceStable(out.Args, func(i, j int) bool {
			return out.Args[i].Type.String() < out.Args[j].Type.String()
		})
		return out
	case "aggregatefunction", "simpleaggregatefunction":
		if len(t.Args) > 0 && t.Args[0].Type != nil && len(t.Args[0].Type.Args) == 0 {
			out := t.Clone()
			name := out.Args[0].Type.Name
			out.Args[0] = core.TypeArg{Ident: &name}
			return out
		}
	}
	return t
}

// resultType is what a conversion returns when that depends on an
// argument's value: toDecimal64(x, s) is Decimal(18, s) and
// toDateTime64(x, p) is DateTime64(p).
func resultType(name string, args []core.ResultArg) *core.TypeExpr {
	switch strings.ToLower(name) {
	case "todecimal32", "todecimal64", "todecimal128", "todecimal256":
		if len(args) >= 2 && args[1].Int != nil {
			p := decimalPrecisions[strings.TrimPrefix(strings.ToLower(name), "to")]
			s := *args[1].Int
			return &core.TypeExpr{Name: "decimal", Args: []core.TypeArg{{Int: &p}, {Int: &s}}}
		}
	case "todatetime64":
		if len(args) >= 2 && args[1].Int != nil {
			p := *args[1].Int
			return &core.TypeExpr{Name: "datetime64", Args: []core.TypeArg{{Int: &p}}}
		}
	}
	return nil
}
