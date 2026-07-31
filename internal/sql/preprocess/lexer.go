package preprocess

import (
	"strings"
)

// lexer is a dialect-parameterized scanner over a SQL file. It does not parse
// SQL; it knows just enough about comments, string literals, quoted identifiers
// and dollar-quoted strings to walk a query without ever looking inside a
// region where sqlc syntax must be left alone.
type lexer struct {
	d   Dialect
	src string
}

func isIdentStart(c byte) bool {
	return c == '_' ||
		(c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c >= 0x80
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || (c >= '0' && c <= '9') || c == '$'
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\f' || c == '\v'
}

func (l *lexer) at(i int) byte {
	if i < 0 || i >= len(l.src) {
		return 0
	}
	return l.src[i]
}

// skip returns the offset just past the comment, string literal, quoted
// identifier or dollar-quoted string that starts at i. It returns -1 when no
// such token starts at i.
func (l *lexer) skip(i int) int {
	switch c := l.at(i); c {
	case '\'':
		return l.skipQuoted(i, '\'', l.d.Backslash || l.isEscapeString(i))
	case '"':
		return l.skipQuoted(i, '"', l.d.DoubleQuoteString && l.d.Backslash)
	case '`':
		if !l.d.Backtick {
			return -1
		}
		return l.skipQuoted(i, '`', false)
	case '-':
		if l.at(i+1) == '-' {
			return l.skipLine(i)
		}
	case '#':
		if l.d.HashComment {
			return l.skipLine(i)
		}
	case '/':
		if l.at(i+1) == '*' {
			return l.skipBlockComment(i)
		}
	case '$':
		if l.d.DollarQuote {
			if end := l.skipDollarQuoted(i); end >= 0 {
				return end
			}
		}
	}
	return -1
}

// isEscapeString reports whether the string literal starting at i is a
// PostgreSQL escape string constant, written E'...', where a backslash escapes
// the following character.
func (l *lexer) isEscapeString(i int) bool {
	if !l.d.DollarQuote { // PostgreSQL-only syntax
		return false
	}
	c := l.at(i - 1)
	if c != 'e' && c != 'E' {
		return false
	}
	// The E must not be the tail of a longer identifier.
	return !isIdentPart(l.at(i - 2))
}

func (l *lexer) skipQuoted(i int, quote byte, backslash bool) int {
	for j := i + 1; j < len(l.src); j++ {
		switch l.src[j] {
		case '\\':
			if backslash {
				j++
			}
		case quote:
			// A doubled quote is an escaped quote, not the end.
			if l.at(j+1) == quote {
				j++
				continue
			}
			return j + 1
		}
	}
	return len(l.src)
}

func (l *lexer) skipLine(i int) int {
	if idx := strings.IndexByte(l.src[i:], '\n'); idx >= 0 {
		return i + idx + 1
	}
	return len(l.src)
}

func (l *lexer) skipBlockComment(i int) int {
	depth := 0
	for j := i; j < len(l.src)-1; j++ {
		if l.src[j] == '/' && l.src[j+1] == '*' {
			depth++
			j++
			continue
		}
		if l.src[j] == '*' && l.src[j+1] == '/' {
			depth--
			j++
			if depth == 0 || !l.d.NestedBlockComment {
				return j + 1
			}
		}
	}
	return len(l.src)
}

// skipDollarQuoted handles PostgreSQL's $tag$ ... $tag$ string literals. It
// returns -1 when i does not start a dollar quote.
func (l *lexer) skipDollarQuoted(i int) int {
	tag := l.dollarTag(i)
	if tag == "" {
		return -1
	}
	if idx := strings.Index(l.src[i+len(tag):], tag); idx >= 0 {
		return i + len(tag) + idx + len(tag)
	}
	return len(l.src)
}

// dollarTag returns the full opening delimiter ("$$" or "$tag$") when one
// starts at i, and "" otherwise.
func (l *lexer) dollarTag(i int) string {
	if l.at(i) != '$' {
		return ""
	}
	j := i + 1
	for j < len(l.src) && l.src[j] != '$' {
		c := l.src[j]
		if isDigit(c) && j == i+1 {
			return "" // $1 is a bind parameter, not a tag
		}
		if !isIdentPart(c) || c == '$' {
			return ""
		}
		j++
	}
	if j >= len(l.src) {
		return ""
	}
	return l.src[i : j+1]
}

// identEnd returns the offset just past the identifier that starts at i.
func (l *lexer) identEnd(i int) int {
	j := i
	for j < len(l.src) && isIdentPart(l.src[j]) {
		j++
	}
	return j
}

// skipSpace advances past whitespace and comments.
func (l *lexer) skipSpace(i int) int {
	for i < len(l.src) {
		if isSpace(l.src[i]) {
			i++
			continue
		}
		if l.src[i] == '-' && l.at(i+1) == '-' {
			i = l.skipLine(i)
			continue
		}
		if l.src[i] == '/' && l.at(i+1) == '*' {
			i = l.skipBlockComment(i)
			continue
		}
		if l.src[i] == '#' && l.d.HashComment {
			i = l.skipLine(i)
			continue
		}
		break
	}
	return i
}

// statements splits the source into statement spans at top-level semicolons.
// The spans cover the whole source, so every byte belongs to exactly one
// statement.
func (l *lexer) statements() [][2]int {
	var out [][2]int
	start := 0
	for i := 0; i < len(l.src); {
		if end := l.skip(i); end >= 0 {
			i = end
			continue
		}
		if l.src[i] == ';' {
			out = append(out, [2]int{start, i + 1})
			start = i + 1
		}
		i++
	}
	if start < len(l.src) {
		out = append(out, [2]int{start, len(l.src)})
	}
	return out
}

// matchParen returns the offset of the ")" that closes the "(" at i, or -1.
func (l *lexer) matchParen(i int) int {
	depth := 0
	for j := i; j < len(l.src); {
		if end := l.skip(j); end >= 0 {
			j = end
			continue
		}
		switch l.src[j] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return j
			}
		}
		j++
	}
	return -1
}

// splitArgs splits the text between parentheses on top-level commas. Empty or
// whitespace-only input yields no arguments.
func (l *lexer) splitArgs(start, end int) []string {
	var args []string
	depth := 0
	last := start
	for j := start; j < end; {
		if stop := l.skip(j); stop >= 0 {
			if stop > end {
				stop = end
			}
			j = stop
			continue
		}
		switch l.src[j] {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				args = append(args, l.src[last:j])
				last = j + 1
			}
		}
		j++
	}
	tail := l.src[last:end]
	if len(args) == 0 && strings.TrimSpace(tail) == "" {
		return nil
	}
	return append(args, tail)
}
