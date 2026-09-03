package core

import (
	"strconv"
	"strings"
)

// TypeExpr is a type written as a call expression: a lowercased name applied
// to arguments that are other types, integers, booleans or strings, each
// with an optional label, and a nullable flag at whatever depth it applies.
// Nothing about a nested type is special-cased, so an array of nullable
// strings is array(string nullable), a map is map(string, uint32) and a
// named tuple is tuple(lat: float64, lon: float64). The catalog resolves the
// names; the expression only records what was declared or inferred.
type TypeExpr struct {
	Name     string    `json:"name"`
	Nullable bool      `json:"nullable,omitempty"`
	Args     []TypeArg `json:"args,omitempty"`
}

// TypeArg is one argument of a TypeExpr: exactly one of Type, Int, Bool or
// String is set.
type TypeArg struct {
	Label  string    `json:"label,omitempty"`
	Type   *TypeExpr `json:"type,omitempty"`
	Int    *int64    `json:"int,omitempty"`
	Bool   *bool     `json:"bool,omitempty"`
	String *string   `json:"string,omitempty"`
}

// ParseTypeExpr reads a type spelled the way every dialect spells one, as a
// name optionally applied to a parenthesised, comma-separated argument list
// whose entries may be labelled (`lat Float64` in a tuple, `'a' = 1` in an
// enum). A Nullable(T) wrapper becomes T with Nullable set, and a trailing
// [] becomes an array of the element.
func ParseTypeExpr(s string) *TypeExpr {
	s = strings.TrimSpace(s)
	if element, ok := strings.CutSuffix(s, ArraySuffix); ok {
		return &TypeExpr{Name: "array", Args: []TypeArg{{Type: ParseTypeExpr(element)}}}
	}
	name, args := splitTypeArgs(s)
	name = strings.ToLower(name)
	if name == "nullable" && len(args) == 1 {
		t := ParseTypeExpr(args[0])
		t.Nullable = true
		return t
	}
	t := &TypeExpr{Name: name}
	for _, a := range args {
		t.Args = append(t.Args, parseTypeArg(a))
	}
	return t
}

func parseTypeArg(a string) TypeArg {
	a = strings.TrimSpace(a)
	if strings.HasPrefix(a, "'") {
		end := quotedEnd(a)
		lit := strings.ReplaceAll(strings.ReplaceAll(a[1:end-1], `\'`, `'`), `''`, `'`)
		if rest := strings.TrimSpace(a[end:]); strings.HasPrefix(rest, "=") {
			arg := parseTypeArg(rest[1:])
			arg.Label = lit
			return arg
		}
		return TypeArg{String: &lit}
	}
	if n, err := strconv.ParseInt(a, 10, 64); err == nil {
		return TypeArg{Int: &n}
	}
	switch strings.ToLower(a) {
	case "true", "false":
		b := strings.EqualFold(a, "true")
		return TypeArg{Bool: &b}
	}
	// A label is a word before a space that comes before any parenthesis,
	// as in `lat Float64` or `tags Array(String)`.
	head := a
	if p := strings.IndexByte(a, '('); p >= 0 {
		head = a[:p]
	}
	if i := strings.IndexByte(head, ' '); i > 0 {
		arg := parseTypeArg(a[i+1:])
		arg.Label = a[:i]
		return arg
	}
	return TypeArg{Type: ParseTypeExpr(a)}
}

// quotedEnd returns the index just past the single-quoted literal starting
// at the beginning of s, honouring backslash escapes and doubled quotes.
func quotedEnd(s string) int {
	for i := 1; i < len(s); i++ {
		switch {
		case s[i] == '\\' && i+1 < len(s):
			i++
		case s[i] == '\'' && i+1 < len(s) && s[i+1] == '\'':
			i++
		case s[i] == '\'':
			return i + 1
		}
	}
	return len(s)
}

// splitTypeArgs splits `Base(arg, arg)` into its base name and top-level
// arguments, leaving nested parentheses and quoted strings intact.
func splitTypeArgs(t string) (string, []string) {
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
	if last := strings.TrimSpace(inner[start:]); last != "" || len(args) > 0 {
		args = append(args, last)
	}
	return base, args
}

// HasNullable reports whether the expression marks nullability anywhere,
// which tells whether the spelling it came from said so itself.
func (t *TypeExpr) HasNullable() bool {
	if t == nil {
		return false
	}
	if t.Nullable {
		return true
	}
	for _, a := range t.Args {
		if a.Type.HasNullable() {
			return true
		}
	}
	return false
}

// String renders the expression in its canonical text form, with a trailing
// "nullable" marking a nullable type.
func (t *TypeExpr) String() string {
	if t == nil {
		return ""
	}
	var b strings.Builder
	b.WriteString(t.Name)
	if len(t.Args) > 0 {
		b.WriteByte('(')
		for i, a := range t.Args {
			if i > 0 {
				b.WriteString(", ")
			}
			if a.Label != "" {
				b.WriteString(a.Label)
				b.WriteString(": ")
			}
			switch {
			case a.Type != nil:
				b.WriteString(a.Type.String())
			case a.Int != nil:
				b.WriteString(strconv.FormatInt(*a.Int, 10))
			case a.Bool != nil:
				b.WriteString(strconv.FormatBool(*a.Bool))
			case a.String != nil:
				b.WriteString("'" + strings.ReplaceAll(*a.String, "'", `\'`) + "'")
			}
		}
		b.WriteByte(')')
	}
	if t.Nullable {
		b.WriteString(" nullable")
	}
	return b.String()
}
