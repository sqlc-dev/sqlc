package duckdb

import (
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
)

// parseType reads a type the way DuckDB spells one — DECIMAL(10,2),
// INTEGER[], INTEGER[3], STRUCT(a INTEGER, b VARCHAR), MAP(VARCHAR,
// INTEGER), UNION(num INTEGER, str VARCHAR), ENUM('a', 'b'), TIMESTAMP
// WITH TIME ZONE — into an expression in lower case, named the way the
// dialect names it when the spelling is one of the aliases types.jsonl
// lists. DuckDB spells an enum column by its labels whether the schema
// named the type or not, so labels that are those of an enum the schema
// created name that type.
func (a *analyzer) parseType(s string) *analysis.TypeExpr {
	s = strings.TrimSpace(s)
	if strings.HasSuffix(s, "]") {
		open := strings.LastIndexByte(s, '[')
		if open > 0 {
			t := &analysis.TypeExpr{Name: "array", Args: []analysis.TypeArg{{Type: a.parseType(s[:open])}}}
			if n, err := strconv.ParseInt(strings.TrimSpace(s[open+1:len(s)-1]), 10, 64); err == nil {
				t.Args = append(t.Args, analysis.TypeArg{Int: &n})
			}
			return t
		}
	}
	name, args := s, ""
	if open := strings.IndexByte(s, '('); open >= 0 && strings.HasSuffix(s, ")") {
		name, args = strings.TrimSpace(s[:open]), s[open+1:len(s)-1]
	}
	name = strings.ToLower(name)
	if c, ok := a.canonical[name]; ok {
		name = c
	}
	t := &analysis.TypeExpr{Name: name}
	if args == "" {
		return t
	}
	switch name {
	case "struct", "union":
		for _, f := range splitTop(args, ',') {
			f = strings.TrimSpace(f)
			label, typ := f, ""
			if strings.HasPrefix(f, `"`) {
				end := quotedEnd(f, 0)
				label, typ = strings.ReplaceAll(f[1:end-1], `""`, `"`), f[end:]
			} else if i := strings.IndexByte(f, ' '); i > 0 {
				label, typ = f[:i], f[i+1:]
			}
			t.Args = append(t.Args, analysis.TypeArg{Label: label, Type: a.parseType(typ)})
		}
	case "enum":
		var labels []string
		for _, l := range splitTop(args, ',') {
			l = strings.TrimSpace(l)
			if strings.HasPrefix(l, "'") && strings.HasSuffix(l, "'") && len(l) >= 2 {
				l = strings.ReplaceAll(l[1:len(l)-1], "''", "'")
			}
			labels = append(labels, l)
		}
		if named, ok := a.enumNamed(labels); ok {
			return &analysis.TypeExpr{Name: named}
		}
		for _, l := range labels {
			l := l
			t.Args = append(t.Args, analysis.TypeArg{String: &l})
		}
	default:
		for _, arg := range splitTop(args, ',') {
			arg = strings.TrimSpace(arg)
			if n, err := strconv.ParseInt(arg, 10, 64); err == nil {
				t.Args = append(t.Args, analysis.TypeArg{Int: &n})
			} else {
				t.Args = append(t.Args, analysis.TypeArg{Type: a.parseType(arg)})
			}
		}
	}
	return t
}

// enumNamed finds the enum type the schema created with these labels.
func (a *analyzer) enumNamed(labels []string) (string, bool) {
	for name, have := range a.enums {
		if len(have) != len(labels) {
			continue
		}
		same := true
		for i := range have {
			if have[i] != labels[i] {
				same = false
				break
			}
		}
		if same {
			return name, true
		}
	}
	return "", false
}

// splitTop splits on a separator outside parentheses, brackets and
// quotes.
func splitTop(s string, sep byte) []string {
	var out []string
	depth, start := 0, 0
	var quote byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case quote != 0:
			if c == quote {
				quote = 0
			}
		case c == '\'' || c == '"':
			quote = c
		case c == '(' || c == '[':
			depth++
		case c == ')' || c == ']':
			depth--
		case c == sep && depth == 0:
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	return append(out, s[start:])
}

// zero is a value of a type, spelled as a string DuckDB casts to the type:
// what a parameter is bound to when the query is run to see which of its
// columns can be NULL. A type no value is known for reports false, and
// the parameter is bound to NULL.
func (a *analyzer) zero(t *analysis.TypeExpr) (string, bool) {
	if t == nil {
		return "", false
	}
	switch t.Name {
	case "tinyint", "smallint", "integer", "bigint", "hugeint", "utinyint", "usmallint", "uinteger", "ubigint", "uhugeint",
		"decimal", "float", "double", "bignum", "varint":
		return "1", true
	case "varchar", "text", "char", "bpchar":
		return "x", true
	case "json":
		return "{}", true
	case "blob":
		return "", true
	case "bit":
		return "0", true
	case "boolean":
		return "true", true
	case "uuid":
		return "00000000-0000-0000-0000-000000000000", true
	case "date":
		return "2000-01-01", true
	case "time", "time with time zone":
		return "00:00:00", true
	case "timestamp", "timestamp with time zone", "timestamp_s", "timestamp_ms", "timestamp_ns", "datetime":
		return "2000-01-01 00:00:00", true
	case "interval":
		return "1 day", true
	case "map":
		return "{}", true
	case "array":
		if len(t.Args) == 0 || t.Args[0].Type == nil {
			return "", false
		}
		elem, ok := a.zero(t.Args[0].Type)
		if !ok {
			return "", false
		}
		n := 1
		if len(t.Args) > 1 && t.Args[1].Int != nil {
			n = int(*t.Args[1].Int)
		}
		elems := make([]string, n)
		for i := range elems {
			elems[i] = elem
		}
		return "[" + strings.Join(elems, ", ") + "]", true
	case "struct":
		var fields []string
		for _, f := range t.Args {
			v, ok := a.zero(f.Type)
			if !ok {
				return "", false
			}
			fields = append(fields, f.Label+": "+v)
		}
		return "{" + strings.Join(fields, ", ") + "}", true
	case "union":
		if len(t.Args) == 0 {
			return "", false
		}
		return a.zero(t.Args[0].Type)
	case "enum":
		if len(t.Args) > 0 && t.Args[0].String != nil {
			return *t.Args[0].String, true
		}
	}
	if labels, ok := a.enums[t.Name]; ok && len(labels) > 0 {
		return labels[0], true
	}
	return "", false
}
