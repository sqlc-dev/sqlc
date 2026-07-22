package golang

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

// clickhouseType maps a ClickHouse column type to a Go type. Type names
// arrive lower-cased from the core catalog (e.g. "uint64", "string",
// "datetime"). Nullable columns (NotNull == false) map to a pointer, which
// is how the clickhouse-go driver represents Nullable(T).
func clickhouseType(req *plugin.GenerateRequest, options *opts.Options, col *plugin.Column) string {
	dt := strings.ToLower(sdk.DataType(col.Type))
	notNull := col.NotNull

	switch dt {
	case "uint8":
		return nullable(notNull, "uint8")
	case "uint16":
		return nullable(notNull, "uint16")
	case "uint32":
		return nullable(notNull, "uint32")
	case "uint64":
		return nullable(notNull, "uint64")
	case "int8":
		return nullable(notNull, "int8")
	case "int16":
		return nullable(notNull, "int16")
	case "int32":
		return nullable(notNull, "int32")
	case "int64":
		return nullable(notNull, "int64")
	case "uint128", "uint256", "int128", "int256":
		// Big integers are represented as *big.Int by clickhouse-go.
		return "*big.Int"
	case "float32", "bfloat16":
		return nullable(notNull, "float32")
	case "float64":
		return nullable(notNull, "float64")
	case "bool":
		return nullable(notNull, "bool")
	case "string", "fixedstring":
		return nullable(notNull, "string")
	case "date", "date32", "datetime", "datetime64":
		return nullable(notNull, "time.Time")

	// The following resolve to string for now; richer mappings
	// (decimal.Decimal, uuid.UUID, netip.Addr, json.RawMessage) require
	// wiring their imports into the Go importer and are a follow-up.
	case "decimal", "decimal32", "decimal64", "decimal128", "decimal256",
		"uuid", "ipv4", "ipv6", "json", "enum8", "enum16":
		return nullable(notNull, "string")

	default:
		return "interface{}"
	}
}

// nullable wraps a base Go type in a pointer when the column is nullable.
func nullable(notNull bool, base string) string {
	if notNull {
		return base
	}
	return "*" + base
}
