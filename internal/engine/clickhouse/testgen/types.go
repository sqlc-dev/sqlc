package main

import (
	"strconv"
	"strings"
)

// A type is a call expression, the way ClickHouse itself models one: a
// lowercased name applied to an ordered argument list. Each argument is
// another type, an integer, a boolean or a quoted string, optionally
// labelled, so
// Nullable, Array and LowCardinality are ordinary names and nothing about a
// nested type is lost. The shape maps one to one onto a protobuf message
// with a oneof for the argument value:
//
//	Map(String, Nullable(UInt32))
//	{"name": "map", "args": [
//	  {"type": {"name": "string"}},
//	  {"type": {"name": "nullable", "args": [{"type": {"name": "uint32"}}]}}]}
//
//	Tuple(lat Float64, lon Float64)
//	{"name": "tuple", "args": [
//	  {"label": "lat", "type": {"name": "float64"}},
//	  {"label": "lon", "type": {"name": "float64"}}]}
//
//	Enum8('a' = 1, 'b' = 2)
//	{"name": "enum8", "args": [{"label": "a", "int": 1}, {"label": "b", "int": 2}]}
//
//	DateTime64(3, 'UTC')
//	{"name": "datetime64", "args": [{"int": 3}, {"string": "UTC"}]}
//
// An identifier argument such as the function in AggregateFunction(uniq,
// String) is a type with no arguments. Resolving names against the catalog
// is the reader's job; the output only records what was said.

type typeExpr struct {
	Name string    `json:"name"`
	Args []typeArg `json:"args,omitempty"`
}

type typeArg struct {
	Label  string    `json:"label,omitempty"`
	Type   *typeExpr `json:"type,omitempty"`
	Int    *int64    `json:"int,omitempty"`
	Bool   *bool     `json:"bool,omitempty"`
	String *string   `json:"string,omitempty"`
}

// parseType turns a ClickHouse type string into its expression.
func parseType(t string) *typeExpr {
	name, args := splitType(t)
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		name = "nothing"
	}
	expr := &typeExpr{Name: name}
	for _, a := range args {
		expr.Args = append(expr.Args, parseArg(a))
	}
	return expr
}

// parseArg parses one argument: a quoted string, an integer, a boolean, a
// labelled argument (`lat Float64` in a Tuple, `'a' = 1` in an Enum), or a
// type.
func parseArg(a string) typeArg {
	a = strings.TrimSpace(a)
	if strings.HasPrefix(a, "'") {
		end := skipQuoted(a, 0)
		lit := unquote(a[1 : end-1])
		if rest := strings.TrimSpace(a[end:]); strings.HasPrefix(rest, "=") {
			arg := parseArg(rest[1:])
			arg.Label = lit
			return arg
		}
		return typeArg{String: &lit}
	}
	if n, err := strconv.ParseInt(a, 10, 64); err == nil {
		return typeArg{Int: &n}
	}
	switch strings.ToLower(a) {
	case "true", "false":
		b := strings.EqualFold(a, "true")
		return typeArg{Bool: &b}
	}
	if i := labelEnd(a); i > 0 {
		arg := parseArg(a[i+1:])
		arg.Label = a[:i]
		return arg
	}
	return typeArg{Type: parseType(a)}
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

// unquote undoes the escaping inside a single-quoted ClickHouse literal.
func unquote(s string) string {
	s = strings.ReplaceAll(s, `\'`, `'`)
	s = strings.ReplaceAll(s, `''`, `'`)
	return strings.ReplaceAll(s, `\\`, `\`)
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
