// Package lex scans SQL source for the trivia the engine parsers discard —
// today, comments. It is the sqlc-side stand-in for the parsers' own token
// streams: as each owned parser learns to expose the comments it already
// lexes, its engine switches from this scanner to the parser's stream and
// nothing downstream changes.
package lex

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// Dialect configures the lexical rules that decide what is a comment and
// what is quoted content.
type Dialect struct {
	// Hash marks # as a line comment (MySQL).
	Hash bool
	// Backslash escapes the next character inside single-quoted strings.
	Backslash bool
	// DollarQuotes recognizes PostgreSQL $tag$ ... $tag$ strings.
	DollarQuotes bool
	// NestedBlocks nests /* */ comments (PostgreSQL).
	NestedBlocks bool
	// Brackets quotes identifiers as [name] (SQLite, T-SQL).
	Brackets bool
}

// SQLite returns the lexical rules of SQLite's tokenizer: no hash comments,
// no backslash escapes, no dollar quotes, no nested block comments,
// [bracketed] identifiers allowed.
func SQLite() Dialect {
	return Dialect{Brackets: true}
}

// Comments returns every comment in src, in order. Offsets are byte offsets
// into src. An unterminated block comment runs to the end of the input,
// matching SQLite and PostgreSQL.
func Comments(d Dialect, src string) []ast.Comment {
	var out []ast.Comment
	n := len(src)
	i := 0
	for i < n {
		c := src[i]
		switch {
		case c == '-' && i+1 < n && src[i+1] == '-':
			start := i
			for i < n && src[i] != '\n' {
				i++
			}
			out = append(out, comment(src, start, i))
		case c == '#' && d.Hash:
			start := i
			for i < n && src[i] != '\n' {
				i++
			}
			out = append(out, comment(src, start, i))
		case c == '/' && i+1 < n && src[i+1] == '*':
			start := i
			i += 2
			depth := 1
			for i < n && depth > 0 {
				switch {
				case src[i] == '*' && i+1 < n && src[i+1] == '/':
					depth--
					i += 2
				case d.NestedBlocks && src[i] == '/' && i+1 < n && src[i+1] == '*':
					depth++
					i += 2
				default:
					i++
				}
			}
			out = append(out, comment(src, start, i))
		case c == '\'' || c == '"' || c == '`':
			i = skipQuoted(src, i, c, d.Backslash && c == '\'')
		case c == '[' && d.Brackets:
			i++
			for i < n && src[i] != ']' {
				i++
			}
			if i < n {
				i++
			}
		case c == '$' && d.DollarQuotes:
			j := i + 1
			for j < n && isTagChar(src[j]) {
				j++
			}
			if j < n && src[j] == '$' {
				tag := src[i : j+1]
				end := strings.Index(src[j+1:], tag)
				if end < 0 {
					return out // unterminated: nothing after is classifiable
				}
				i = j + 1 + end + len(tag)
			} else {
				i++
			}
		default:
			i++
		}
	}
	return out
}

func comment(src string, start, end int) ast.Comment {
	// A comment is on its own line when nothing but blank space sits
	// between the preceding line break and the comment.
	ownLine := true
	for j := start - 1; j >= 0; j-- {
		ch := src[j]
		if ch == '\n' {
			break
		}
		if ch != ' ' && ch != '\t' && ch != '\r' {
			ownLine = false
			break
		}
	}
	return ast.Comment{
		Text:    strings.TrimRight(src[start:end], " \t\r"),
		Start:   start,
		End:     end,
		OwnLine: ownLine,
	}
}

func skipQuoted(src string, i int, term byte, backslash bool) int {
	n := len(src)
	i++
	for i < n {
		if backslash && src[i] == '\\' {
			i += 2
			continue
		}
		if src[i] == term {
			// A doubled quote is an escaped quote, not the end.
			if i+1 < n && src[i+1] == term {
				i += 2
				continue
			}
			return i + 1
		}
		i++
	}
	return i
}

func isTagChar(c byte) bool {
	return c == '_' ||
		('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z') ||
		('0' <= c && c <= '9')
}
