package mysql

import (
	"fmt"
)

// converts MySQL function name to MySQL return type
func functionReturnType(f string) string {
	switch f {
	case "avg", "count", "instr", "sum", "min", "max", "length", "char_length":
		return "int"
	case "concat", "left", "replace", "substring", "trim", "find_in_set", "format", "group_concat":
		return "varchar"
	default:
		panic(fmt.Sprintf("unknown mysql function type [%v]", f))
	}
}
