package main

import "strings"

// A type is written as a call expression, the way ClickHouse itself models
// one: a lowercased name applied to an ordered argument list in which each
// argument is a number, a quoted string, an identifier, another call, or
// one of those with a label. Nothing is special-cased, so Nullable, Array
// and LowCardinality are ordinary names, and the text form is the same for
// every engine that spells the same structure differently:
//
//	Map(String, Nullable(UInt32))     map(string, nullable(uint32))
//	Tuple(lat Float64, lon Float64)   tuple(lat: float64, lon: float64)
//	Enum8('a' = 1, 'b' = 2)           enum8('a': 1, 'b': 2)
//	DateTime64(3, 'UTC')              datetime64(3, 'UTC')
//	AggregateFunction(uniq, String)   aggregatefunction(uniq, string)
//
// The catalog resolves the names afterwards; the output only records what
// was said.

// typeExpr renders a column's type and reports whether it is NOT NULL. An
// outer Nullable is the column's nullability rather than part of its type,
// so it is lifted into the flag, through LowCardinality when needed; a
// Nullable anywhere deeper stays in the expression.
func typeExpr(t string) (expr string, notNull bool) {
	name, args := splitType(t)
	switch strings.ToLower(name) {
	case "nullable":
		if len(args) == 1 {
			expr, _ = typeExpr(args[0])
			return expr, false
		}
	case "lowcardinality":
		if len(args) == 1 {
			inner, notNull := typeExpr(args[0])
			return "lowcardinality(" + inner + ")", notNull
		}
	}
	return renderCall(t), true
}

// renderCall renders a type as a call expression.
func renderCall(t string) string {
	name, args := splitType(t)
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return "nothing"
	}
	if args == nil {
		return name
	}
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = renderArg(a)
	}
	return name + "(" + strings.Join(parts, ", ") + ")"
}

// renderArg renders one argument: a quoted string, a number, a labelled
// argument (`lat Float64` in a Tuple, `'a' = 1` in an Enum), or a call.
func renderArg(a string) string {
	a = strings.TrimSpace(a)
	if strings.HasPrefix(a, "'") {
		end := skipQuoted(a, 0)
		if rest := strings.TrimSpace(a[end:]); strings.HasPrefix(rest, "=") {
			return a[:end] + ": " + renderArg(rest[1:])
		}
		return a[:end]
	}
	if isNumber(a) {
		return a
	}
	if i := labelEnd(a); i > 0 {
		return a[:i] + ": " + renderArg(a[i+1:])
	}
	return renderCall(a)
}

// labelEnd returns the index of the space separating a label from the type
// it labels, or -1 when the argument has no label: a space that comes before
// any parenthesis, as in `lat Float64` or `tags Array(String)`.
func labelEnd(a string) int {
	head := a
	if p := strings.IndexByte(a, '('); p >= 0 {
		head = a[:p]
	}
	return strings.IndexByte(head, ' ')
}

func isNumber(s string) bool {
	if s == "" {
		return false
	}
	for i, c := range s {
		if c == '-' && i == 0 && len(s) > 1 {
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
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
