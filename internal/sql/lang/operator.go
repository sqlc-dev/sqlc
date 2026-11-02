package lang

// TODO: This logic is PostgreSQL-specific and needs to be refactored to support MySQL

func IsComparisonOperator(s string) bool {
	switch s {
	case ">":
	case "<":
	case "<=":
	case ">=":
	case "=":
	case "==":
	case "<>":
	case "!=":
	// Comparisons an engine spells as a keyword rather than as punctuation.
	// They are written upper case, which is how SQLite's parser records
	// them; PostgreSQL and MySQL use distinct node types or their own
	// operator names for these, so the spellings do not overlap.
	case "IS":
	case "IS NOT":
	case "IS DISTINCT FROM":
	case "IS NOT DISTINCT FROM":
	case "ISNULL":
	case "NOTNULL":
	case "NOT NULL":
	case "LIKE":
	case "NOT LIKE":
	case "GLOB":
	case "NOT GLOB":
	case "REGEXP":
	case "NOT REGEXP":
	case "MATCH":
	case "NOT MATCH":
	default:
		return false
	}
	return true
}

func IsMathematicalOperator(s string) bool {
	switch s {
	case "+":
	case "-":
	case "*":
	case "/":
	case "%":
	case "^":
	case "|/":
	case "||/":
	case "!":
	case "!!":
	case "@":
	case "&":
	case "|":
	case "#":
	case "~":
	case "<<":
	case ">>":
	default:
		return false
	}
	return true
}

// IsJSONNullableOperator reports whether op is a Postgres JSON / JSONB
// accessor operator that returns SQL NULL when the requested key or path
// is missing from the value. These are: -> (object/array element as jsonb),
// ->> (object/array element as text), #> (path lookup as jsonb), and #>>
// (path lookup as text). Wrapping such an expression in a type cast does
// not make the result non-nullable — the key may still be absent at run
// time. See issue #3792.
func IsJSONNullableOperator(s string) bool {
	switch s {
	case "->", "->>", "#>", "#>>":
		return true
	}
	return false
}
