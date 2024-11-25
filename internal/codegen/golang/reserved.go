package golang

func escape(s string) string {
	if IsReserved(s) {
		return s + "_"
	}
	return s
}

func IsReserved(s string) bool {
	switch s {
	case "break":
		return true
	case "default":
		return true
	case "func":
		return true
	case "interface":
		return true
	case "select":
		return true
	case "case":
		return true
	case "defer":
		return true
	case "go":
		return true
	case "map":
		return true
	case "struct":
		return true
	case "chan":
		return true
	case "else":
		return true
	case "goto":
		return true
	case "package":
		return true
	case "switch":
		return true
	case "const":
		return true
	case "fallthrough":
		return true
	case "if":
		return true
	case "range":
		return true
	case "type":
		return true
	case "continue":
		return true
	case "for":
		return true
	case "import":
		return true
	case "return":
		return true
	case "var":
		return true
	case "q":
		return true
	default:
		return false
	}
}
