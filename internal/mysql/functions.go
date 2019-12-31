package mysql

func functionReturnType(f string) string {
	switch f {
	case "avg", "count", "instr", "sum", "min", "max", "length", "char_length":
		return "int"
	case "concat", "left", "replace", "substring", "trim", "find_in_set", "format":
		return "varchar"
	default:
		panic("unknown mysql function type")
	}
}
