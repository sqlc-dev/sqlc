package preprocess

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

// scan walks a single statement and returns every sqlc construct and native
// placeholder it contains, in source order.
func (l *lexer) scan(start, end int) []occurrence {
	var out []occurrence
	for i := start; i < end; {
		if stop := l.skip(i); stop >= 0 {
			if stop <= i {
				stop = i + 1
			}
			i = stop
			continue
		}
		c := l.src[i]
		switch {
		case isIdentStart(c):
			j := l.identEnd(i)
			if strings.EqualFold(l.src[i:j], "sqlc") {
				if occ, next, ok := l.sqlcCall(i, j); ok {
					out = append(out, occ)
					i = next
					continue
				}
			}
			i = j

		case c == '@' && l.d.AtSign && isIdentStart(l.at(i+1)):
			j := l.identEnd(i + 1)
			out = append(out, occurrence{
				kind:  kindArg,
				start: i,
				end:   j,
				name:  l.src[i+1 : j],
				ident: l.src[i:j],
			})
			i = j

		case c == '$' && l.d.DollarNumber && isDigit(l.at(i+1)):
			j := i + 1
			for j < len(l.src) && isDigit(l.src[j]) {
				j++
			}
			n, _ := strconv.Atoi(l.src[i+1 : j])
			out = append(out, occurrence{
				kind:     kindNative,
				start:    i,
				end:      j,
				number:   n,
				explicit: true,
				ident:    l.src[i:j],
			})
			i = j

		case c == '?' && l.d.Question:
			j := i + 1
			for j < len(l.src) && isDigit(l.src[j]) {
				j++
			}
			occ := occurrence{kind: kindNative, start: i, end: j, ident: l.src[i:j]}
			if j > i+1 {
				occ.number, _ = strconv.Atoi(l.src[i+1 : j])
				occ.explicit = true
			}
			out = append(out, occ)
			i = j

		default:
			i++
		}
	}
	return out
}

// sqlcCall parses a "sqlc.<name>(<arg>)" call that starts at start, where
// afterSqlc is the offset just past the "sqlc" identifier. It reports false
// when the text is not a sqlc call at all, in which case scanning continues
// normally.
func (l *lexer) sqlcCall(start, afterSqlc int) (occurrence, int, bool) {
	i := l.skipSpace(afterSqlc)
	if l.at(i) != '.' {
		return occurrence{}, 0, false
	}
	i = l.skipSpace(i + 1)
	if !isIdentStart(l.at(i)) {
		return occurrence{}, 0, false
	}
	nameEnd := l.identEnd(i)
	fn := strings.ToLower(l.src[i:nameEnd])

	open := l.skipSpace(nameEnd)
	if l.at(open) != '(' {
		return occurrence{}, 0, false
	}
	closing := l.matchParen(open)
	if closing < 0 {
		return occurrence{}, 0, false
	}

	occ := occurrence{start: start, end: closing + 1}
	switch fn {
	case "arg":
		occ.kind = kindArg
	case "narg":
		occ.kind = kindNarg
	case "slice":
		occ.kind = kindSlice
	case "embed":
		occ.kind = kindEmbed
	default:
		occ.err = sqlerr.FunctionNotFound("sqlc." + fn)
		occ.err.Location = start
		return occ, occ.end, true
	}

	args := l.splitArgs(open+1, closing)
	if len(args) != 1 {
		occ.err = &sqlerr.Error{
			Message:  fmt.Sprintf("expected 1 parameter to sqlc.%s; got %d", fn, len(args)),
			Location: start,
		}
		return occ, occ.end, true
	}

	name, ident, ok := argument(l.d, args[0])
	if !ok {
		occ.err = &sqlerr.Error{
			Message:  fmt.Sprintf("expected parameter to sqlc.%s to be string or reference; got %q", fn, strings.TrimSpace(args[0])),
			Location: start,
		}
		return occ, occ.end, true
	}
	occ.name = name
	occ.ident = ident
	return occ, occ.end, true
}

// argument interprets the single argument to a sqlc function. It returns the
// parameter name and the text to reuse when rebuilding a reference, which keeps
// any identifier quoting the user wrote.
func argument(d Dialect, arg string) (name, ident string, ok bool) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return "", "", false
	}
	switch arg[0] {
	case '\'':
		// A string constant names the parameter but is not a reference, so the
		// rebuilt identifier drops the quotes.
		s, ok := unquote(arg, '\'')
		return s, s, ok
	case '"':
		s, ok := unquote(arg, '"')
		if !ok {
			return "", "", false
		}
		if d.DoubleQuoteString {
			return s, s, true
		}
		return s, arg, true
	case '`':
		if !d.Backtick {
			return "", "", false
		}
		s, ok := unquote(arg, '`')
		if !ok {
			return "", "", false
		}
		return s, arg, true
	}
	if !isIdentStart(arg[0]) {
		return "", "", false
	}
	for i := 0; i < len(arg); i++ {
		if !isIdentPart(arg[i]) && arg[i] != '.' {
			return "", "", false
		}
	}
	// A qualified reference flattens to a single name, matching how the AST
	// rewrite joined the parts of a ColumnRef.
	name = strings.ReplaceAll(arg, ".", "")
	if d.FoldIdentifier {
		name = strings.ToLower(name)
	}
	return name, arg, true
}

func unquote(s string, quote byte) (string, bool) {
	if len(s) < 2 || s[0] != quote || s[len(s)-1] != quote {
		return "", false
	}
	inner := s[1 : len(s)-1]
	if strings.IndexByte(inner, quote) >= 0 {
		doubled := string([]byte{quote, quote})
		if strings.Contains(inner, doubled) {
			return strings.ReplaceAll(inner, doubled, string(quote)), true
		}
		return "", false
	}
	return inner, true
}
