package postgres

func args(a ...int) map[int]struct{} {
	m := map[int]struct{}{}
	for _, arg := range a {
		m[arg] = struct{}{}
	}
	return m
}

var Functions = map[string]map[int]struct{}{

	// https://www.postgresql.org/docs/current/functions-math.html
	// Table 9.5. Mathematical Functions
	"abs":          args(1),
	"cbrt":         args(1),
	"ceil":         args(1),
	"ceiling":      args(1),
	"degrees":      args(1),
	"div":          args(2),
	"exp":          args(1),
	"floor":        args(1),
	"ln":           args(1),
	"log":          args(1, 2),
	"mod":          args(2),
	"pi":           args(0),
	"power":        args(2),
	"radians":      args(1),
	"round":        args(1, 2),
	"scale":        args(1),
	"sign":         args(1),
	"sqrt":         args(1),
	"trunc":        args(1, 2),
	"width_bucket": args(2, 4),

	// Table 9.6. Random Functions
	"random":  args(0),
	"setseed": args(1),

	// Table 9.7. Trigonometric Functions
	"acos":   args(1),
	"acosd":  args(1),
	"asin":   args(1),
	"asind":  args(1),
	"atan":   args(1),
	"atan2":  args(2),
	"atan2d": args(2),
	"atand":  args(1),
	"cos":    args(1),
	"cosd":   args(1),
	"cot":    args(1),
	"cotd":   args(1),
	"sin":    args(1),
	"sind":   args(1),
	"tan":    args(1),
	"tand":   args(1),

	// https://www.postgresql.org/docs/current/functions-string.html
	// Table 9.8. SQL String Functions and Operators
	"bit_length":          args(1),
	"char_length":         args(1),
	"character_length":    args(1),
	"lower":               args(1),
	"octet_length":        args(1),
	"overlay":             args(3, 4),
	"pg_catalog.position": args(2),
	"substring":           args(1, 2, 3),
	"trim":                args(2, 3),
	"upper":               args(1),

	// Table 9.9. Other String Functions
	"ascii":                 args(1),
	"btrim":                 args(1, 2),
	"chr":                   args(1),
	"convert":               args(3),
	"convert_from":          args(2),
	"convert_to":            args(2),
	"decode":                args(2),
	"encode":                args(2),
	"initcap":               args(1),
	"left":                  args(2),
	"length":                args(1, 2),
	"lpad":                  args(2, 3),
	"ltrim":                 args(1, 2),
	"md5":                   args(1),
	"parse_ident":           args(1, 2),
	"pg_client_encoding":    args(0),
	"quote_ident":           args(1),
	"quote_literal":         args(1),
	"quote_nullable":        args(1),
	"regexp_match":          args(2, 3),
	"regexp_matches":        args(2, 3),
	"regexp_replace":        args(3, 4),
	"regexp_split_to_array": args(2, 3),
	"regexp_split_to_table": args(2, 3),
	"repeat":                args(2),
	"replace":               args(3),
	"reverse":               args(1),
	"right":                 args(2),
	"rpad":                  args(2, 3),
	"rtrim":                 args(1, 2),
	"split_part":            args(3),
	"strpos":                args(2),
	"substr":                args(2, 3),
	"starts_with":           args(2),
	"to_ascii":              args(1, 2),
	"to_hex":                args(1),
	"translate":             args(3),

	// https://www.postgresql.org/docs/current/functions-binarystring.html
	// Table 9.12. Other Binary String Functions
	"get_bit":  args(2),
	"get_byte": args(2),
	"set_bit":  args(3),
	"set_byte": args(3),
	"sha224":   args(1),
	"sha256":   args(1),
	"sha384":   args(1),
	"sha512":   args(1),

	// https://www.postgresql.org/docs/current/functions-formatting.html
	// Table 9.23. Formatting Functions
	"to_char":      args(2),
	"to_date":      args(2),
	"to_number":    args(2),
	"to_timestamp": args(1, 2),

	// https://www.postgresql.org/docs/current/functions-datetime.html
	"age":              args(1, 2),
	"clock_timestamp":  args(0),
	"date_part":        args(2),
	"date_trunc":       args(2),
	"extract":          args(2),
	"isfinite":         args(1),
	"justify_days":     args(1),
	"justify_hours":    args(1),
	"justify_interval": args(1),
	"make_date":        args(3),
	"make_time":        args(3),
	"make_timestamp":   args(6),
	"make_timestampz":  args(6),
	"now":              args(0),
	"statement_timestamp":   args(0),
	"timeofday":             args(0),
	"transaction_timestamp": args(0),

	// https://www.postgresql.org/docs/current/functions-enum.html
	// Table 9.32. Enum Support Functions
	"enum_first": args(1),
	"enum_last":  args(1),
	"enum_range": args(1, 2),

	// https://www.postgresql.org/docs/current/functions-geometry.html
	// Table 9.34. Geometric Functions
	"area":     args(1),
	"center":   args(1),
	"diameter": args(1),
	"height":   args(1),
	"isclosed": args(1),
	"isopen":   args(1),
	"npoints":  args(1),
	"pclose":   args(1),
	"popen":    args(1),
	"radius":   args(1),
	"width":    args(1),

	// Table 9.35. Geometric Type Conversion Functions
	"box":       args(1, 2),
	"bound_box": args(2),
	"circle":    args(1, 2),
	"line":      args(2),
	"lseg":      args(1, 2),
	"path":      args(1),
	"point":     args(1, 2),
	"polygon":   args(1, 2),

	// https://www.postgresql.org/docs/current/functions-net.html
	// Table 9.37. cidr and inet Functions
	"abbrev":           args(1),
	"broadcast":        args(1),
	"family":           args(1),
	"host":             args(1),
	"hostmask":         args(1),
	"masklen":          args(1),
	"netmask":          args(1),
	"network":          args(1),
	"set_masklen":      args(1),
	"text":             args(1),
	"inet_same_family": args(1),
	"inet_merge":       args(1),

	// https://www.postgresql.org/docs/current/functions-aggregate.html
	"count": args(0, 1),
}
