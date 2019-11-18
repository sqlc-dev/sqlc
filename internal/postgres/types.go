package postgres

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
