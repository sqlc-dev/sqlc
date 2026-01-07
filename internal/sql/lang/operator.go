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
	case "div": // TODO: MySQL-specific operator - should be moved to engine-specific logic
	case "mod": // TODO: MySQL-specific operator - should be moved to engine-specific logic
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
