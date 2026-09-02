package main

import "strings"

// sqlcType lowers a ClickHouse type name to the shape sqlc's ClickHouse
// engine reports: the lowercased base name with parameters dropped,
// Nullable and Array wrappers folded into flags and LowCardinality
// discarded. It mirrors unwrapTypeString in the engine so goldens generated
// here diff cleanly against `sqlc analyze`.
func sqlcType(t string) (dataType string, isArray, notNull bool) {
	name, arr, nullable := unwrapType(t)
	return name, arr, !nullable
}

func unwrapType(t string) (name string, isArray, nullable bool) {
	base, args := splitType(t)
	switch strings.ToLower(base) {
	case "nullable":
		if len(args) == 1 {
			inner, arr, _ := unwrapType(args[0])
			return inner, arr, true
		}
		return "nullable", false, true
	case "lowcardinality":
		if len(args) == 1 {
			return unwrapType(args[0])
		}
		return "lowcardinality", false, false
	case "array":
		if len(args) == 1 {
			inner, _, nul := unwrapType(args[0])
			return inner, true, nul
		}
		return "array", true, false
	case "":
		return "nothing", false, false
	default:
		return strings.ToLower(base), false, false
	}
}

// splitType splits `Base(arg, arg)` into its base name and top-level
// arguments, leaving nested parentheses and quoted strings intact.
func splitType(t string) (string, []string) {
	t = strings.TrimSpace(t)
	open := strings.IndexByte(t, '(')
	if open < 0 || !strings.HasSuffix(t, ")") {
		return t, nil
	}
	base := strings.TrimSpace(t[:open])
	inner := t[open+1 : len(t)-1]
	var (
		args  []string
		depth int
		quote byte
		start int
	)
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		switch {
		case quote != 0:
			if c == '\\' {
				i++
			} else if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"' || c == '`':
			quote = c
		case c == '(':
			depth++
		case c == ')':
			depth--
		case c == ',' && depth == 0:
			args = append(args, strings.TrimSpace(inner[start:i]))
			start = i + 1
		}
	}
	args = append(args, strings.TrimSpace(inner[start:]))
	return base, args
}
