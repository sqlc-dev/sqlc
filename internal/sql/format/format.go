package format

// Formatter provides SQL dialect-specific formatting behavior
type Formatter interface {
	// QuoteIdent returns a quoted identifier if it needs quoting
	// (e.g., reserved words, mixed case identifiers)
	QuoteIdent(s string) string

	// TypeName returns the SQL type name for the given namespace and name.
	// This handles dialect-specific type name mappings (e.g., pg_catalog.int4 -> integer)
	TypeName(ns, name string) string
}
