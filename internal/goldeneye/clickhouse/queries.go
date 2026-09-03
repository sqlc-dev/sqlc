package clickhouse

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// query is one entry of a sqlc query file.
type query struct {
	Name string
	Cmd  string
	SQL  string
}

var headerRe = regexp.MustCompile(`^\s*--\s*name:\s*(\S+)\s+(:\S+)\s*$`)

// parseQueries splits a sqlc query file on its `-- name: X :cmd` headers.
func parseQueries(src string) ([]query, error) {
	var queries []query
	var cur *query
	var body []string
	flush := func() error {
		if cur == nil {
			return nil
		}
		sql := strings.TrimSpace(strings.Join(body, "\n"))
		sql = strings.TrimRight(sql, "; \t\r\n")
		if sql == "" {
			return fmt.Errorf("query %s has no body", cur.Name)
		}
		cur.SQL = sql
		queries = append(queries, *cur)
		return nil
	}
	for _, line := range strings.Split(src, "\n") {
		if m := headerRe.FindStringSubmatch(line); m != nil {
			if err := flush(); err != nil {
				return nil, err
			}
			cur = &query{Name: m[1], Cmd: m[2]}
			body = body[:0]
			continue
		}
		if cur != nil {
			body = append(body, line)
		}
	}
	if err := flush(); err != nil {
		return nil, err
	}
	if len(queries) == 0 {
		return nil, fmt.Errorf("no queries found: expected `-- name: Name :cmd` headers")
	}
	return queries, nil
}

// placeholder is one parameter reference in a query, in order of appearance.
type placeholder struct {
	Number int
	Name   string // sqlc.arg / sqlc.narg name, empty for ?
}

// Placeholders are substituted with constant expressions that carry their
// ordinal, so they can be told apart from each other and from literal NULLs
// once ClickHouse has folded them: the query tree prints the expression a
// folded constant came from. NULL coerces to any type, so a comparison
// against it analyzes with the other operand's type. LIMIT and OFFSET
// reject NULL and only accept unsigned integers, so those get a value no
// query would plausibly contain.
const limitBase = "4294967295"

func sentinelFor(lastWord string, ordinal int) string {
	switch strings.ToLower(lastWord) {
	case "limit", "offset":
		return fmt.Sprintf("toUInt64(%s + %d)", limitBase, ordinal)
	}
	return fmt.Sprintf("(NULL + %d)", ordinal)
}

var namedArgRe = regexp.MustCompile(`^sqlc\.(n?arg)\(\s*'?([A-Za-z_][A-Za-z0-9_]*)'?\s*\)`)

// bindPlaceholders rewrites sqlc's parameter syntax (?, sqlc.arg(name),
// sqlc.narg(name)) into constants ClickHouse can analyze, skipping string
// literals, quoted identifiers and comments. ClickHouse binds every ?
// positionally, so each placeholder is its own parameter even when a name
// repeats, which is how sqlc numbers them too.
func bindPlaceholders(sql string) (string, []placeholder) {
	var (
		out      strings.Builder
		phs      []placeholder
		lastWord string
		i        = 0
	)
	number := func(string) int { return len(phs) + 1 }
	for i < len(sql) {
		c := sql[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			end := skipQuoted(sql, i)
			out.WriteString(sql[i:end])
			i = end
		case strings.HasPrefix(sql[i:], "--"):
			end := strings.IndexByte(sql[i:], '\n')
			if end < 0 {
				end = len(sql)
			} else {
				end += i
			}
			out.WriteString(sql[i:end])
			i = end
		case strings.HasPrefix(sql[i:], "/*"):
			end := strings.Index(sql[i:], "*/")
			if end < 0 {
				end = len(sql)
			} else {
				end += i + 2
			}
			out.WriteString(sql[i:end])
			i = end
		case c == '?':
			out.WriteString(sentinelFor(lastWord, len(phs)+1))
			phs = append(phs, placeholder{Number: number("")})
			lastWord = ""
			i++
		case c == 's' && namedArgRe.MatchString(sql[i:]):
			m := namedArgRe.FindStringSubmatch(sql[i:])
			out.WriteString(sentinelFor(lastWord, len(phs)+1))
			phs = append(phs, placeholder{Number: number(m[2]), Name: m[2]})
			lastWord = ""
			i += len(m[0])
		case isWordByte(c):
			end := i
			for end < len(sql) && isWordByte(sql[end]) {
				end++
			}
			lastWord = sql[i:end]
			out.WriteString(lastWord)
			i = end
		default:
			if !unicode.IsSpace(rune(c)) && c != ',' && c != '(' {
				lastWord = ""
			}
			out.WriteByte(c)
			i++
		}
	}
	return out.String(), phs
}

func isWordByte(c byte) bool {
	return c == '_' || c == '.' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// skipQuoted returns the index just past the quoted token starting at i,
// honouring backslash escapes and doubled quotes.
func skipQuoted(s string, i int) int {
	q := s[i]
	j := i + 1
	for j < len(s) {
		switch {
		case s[j] == '\\' && j+1 < len(s):
			j += 2
		case s[j] == q && j+1 < len(s) && s[j+1] == q:
			j += 2
		case s[j] == q:
			return j + 1
		default:
			j++
		}
	}
	return len(s)
}
