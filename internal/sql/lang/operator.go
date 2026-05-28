package lang

// TODO: This logic is PostgreSQL-specific and needs to be refactored to support MySQL

func IsComparisonOperator(s string) bool {
	switch s {
	case ">":
	case "<":
	case "<=":
	case ">=":
	case "=":
	case "<>":
	case "!=":
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
