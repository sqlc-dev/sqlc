package format

// Formatter provides SQL dialect-specific formatting behavior
type Formatter interface {
	// QuoteIdent returns a quoted identifier if it needs quoting
	// (e.g., reserved words, mixed case identifiers)
	QuoteIdent(s string) string
}
