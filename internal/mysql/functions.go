package mysql

import (
	"fmt"
)

// converts MySQL function name to MySQL return type
func functionReturnType(f string) string {
	switch f {
	case "avg", "count", "instr", "sum", "min", "max", "length", "char_length",
		"ceil", "floor", "mod", "isnull":
		return "int"
	case "concat", "left", "replace", "substring", "trim", "find_in_set", "format", "group_concat":
		return "varchar"
	case "abs", "round", "truncate":
		return "decimal"
	default:
		panic(fmt.Sprintf("unknown mysql function type \"%v\"", f))
	}
}

// returns true if MySQL function can return null.
// See: https://dev.mysql.com/doc/refman/8.0/en/function-reference.html
func functionIsNullable(f string) bool {
	switch f {
	case "avg", "sum", "min", "max", "mod",
		"concat", "left", "find_in_set", "group_concat":
		return true
	default:
		return false
	}
}
