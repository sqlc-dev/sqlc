package format

// Dialect provides SQL dialect-specific formatting behavior
type Dialect interface {
	// QuoteIdent returns a quoted identifier if it needs quoting
	// (e.g., reserved words, mixed case identifiers)
	QuoteIdent(s string) string

	// TypeName returns the SQL type name for the given namespace and name.
	// This handles dialect-specific type name mappings (e.g., pg_catalog.int4 -> integer)
	TypeName(ns, name string) string

	// Param returns the placeholder for the parameter with this number.
	// numbered reports that the author wrote the number out (SQLite's ?2,
	// where a bare ? takes the next index), so engines with both forms can
	// keep the one that was written.
	// PostgreSQL uses $1, $2, etc. MySQL uses ?
	Param(n int, numbered bool) string

	// Cast formats a type cast expression.
	// PostgreSQL uses expr::type, MySQL uses CAST(expr AS type)
	Cast(arg, typeName string) string
}
