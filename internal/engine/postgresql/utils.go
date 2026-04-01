package postgresql

import (
	nodes "github.com/pganalyze/pg_query_go/v6"
)

func isArray(n *nodes.TypeName) bool {
	if n == nil {
		return false
	}
	return len(n.ArrayBounds) > 0
}

func isNotNull(n *nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints {
		switch inner := c.Node.(type) {
		case *nodes.Node_Constraint:
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_NOTNULL {
				return true
			}
			if inner.Constraint.Contype == nodes.ConstrType_CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
}

func IsNamedParamFunc(node *nodes.Node) bool {
	fun, ok := node.Node.(*nodes.Node_FuncCall)
	return ok && joinNodes(fun.FuncCall.Funcname, ".") == "sqlc.arg"
}

func IsNamedParamSign(node *nodes.Node) bool {
	expr, ok := node.Node.(*nodes.Node_AExpr)
	return ok && joinNodes(expr.AExpr.Name, ".") == "@"
}

// TypeFamily maps a PostgreSQL DataType string to a canonical type family name,
// grouping compatible type aliases together. This is used for type compatibility
// checks rather than exact string equality, because PostgreSQL considers many
// type aliases assignment-compatible (e.g. text and varchar are both string types).
//
// The groupings are derived from postgresType() in
// internal/codegen/golang/postgresql_type.go, which maps these aliases to the
// same Go type. We cannot call postgresType() directly for type compatibility
// checking because it requires *plugin.GenerateRequest — a protobuf codegen
// struct constructed after compilation — and driver-specific opts.Options.
func TypeFamily(dt string) string {
	switch dt {
	case "serial", "serial4", "pg_catalog.serial4",
		"integer", "int", "int4", "pg_catalog.int4":
		return "int32"
	case "bigserial", "serial8", "pg_catalog.serial8",
		"bigint", "int8", "pg_catalog.int8",
		"interval", "pg_catalog.interval":
		return "int64"
	case "smallserial", "serial2", "pg_catalog.serial2",
		"smallint", "int2", "pg_catalog.int2":
		return "int16"
	case "float", "double precision", "float8", "pg_catalog.float8":
		return "float64"
	case "real", "float4", "pg_catalog.float4":
		return "float32"
	case "numeric", "pg_catalog.numeric", "money":
		return "numeric"
	case "boolean", "bool", "pg_catalog.bool":
		return "bool"
	case "json", "pg_catalog.json":
		return "json"
	case "jsonb", "pg_catalog.jsonb":
		return "jsonb"
	case "bytea", "blob", "pg_catalog.bytea":
		return "bytes"
	case "date":
		return "date"
	case "pg_catalog.time":
		return "time"
	case "pg_catalog.timetz":
		return "timetz"
	case "pg_catalog.timestamp", "timestamp":
		return "timestamp"
	case "pg_catalog.timestamptz", "timestamptz":
		return "timestamptz"
	case "text", "pg_catalog.varchar", "pg_catalog.bpchar",
		"string", "citext", "name",
		"ltree", "lquery", "ltxtquery":
		return "text"
	case "uuid":
		return "uuid"
	case "inet":
		return "inet"
	case "cidr":
		return "cidr"
	case "macaddr", "macaddr8":
		return "macaddr"
	case "bit", "varbit", "pg_catalog.bit", "pg_catalog.varbit":
		return "bits"
	case "hstore":
		return "hstore"
	case "vector":
		return "vector"
	default:
		return dt
	}
}

func makeByte(s string) byte {
	var b byte
	if s == "" {
		return b
	}
	return []byte(s)[0]
}

func makeUint32Slice(in []uint64) []uint32 {
	out := make([]uint32, len(in))
	for i, v := range in {
		out[i] = uint32(v)
	}
	return out
}

func makeString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
