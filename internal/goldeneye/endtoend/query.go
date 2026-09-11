package endtoend

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"
)

// Query is one entry of a sqlc query file.
type Query struct {
	Name string
	Cmd  string
	SQL  string
}

// Queries reads the case's query file and splits it into its queries.
func (c Case) Queries() ([]Query, error) {
	src, err := os.ReadFile(c.Query)
	if err != nil {
		return nil, err
	}
	queries, err := ParseQueries(string(src))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", c.Query, err)
	}
	return queries, nil
}

var headerRe = regexp.MustCompile(`^\s*--\s*name:\s*(\S+)\s+(:\S+)\s*$`)

// ParseQueries splits a sqlc query file on its `-- name: X :cmd` headers.
func ParseQueries(src string) ([]Query, error) {
	var queries []Query
	var cur *Query
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
			cur = &Query{Name: m[1], Cmd: m[2]}
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

var namedArgRe = regexp.MustCompile(`^sqlc\.(n?arg|slice)\(\s*'?([A-Za-z_][A-Za-z0-9_]*)'?\s*\)`)

// typedParamRe matches ClickHouse's {name:Type} parameter, whose type is
// the query's own business: the engine binds it as it binds any other.
var typedParamRe = regexp.MustCompile(`^\{\s*([A-Za-z_][A-Za-z0-9_]*)\s*:\s*[^}]+\}`)

// Rewrite replaces every parameter reference in a query — ?, sqlc.arg(name),
// sqlc.narg(name), sqlc.slice(name) and ClickHouse's {name:Type} — with
// what bind returns for it, in order of appearance, skipping string
// literals, quoted identifiers and comments. bind is handed the name, empty for a ?, and the word before the
// reference, so that a LIMIT or OFFSET can be bound differently from a
// value; the second count of a LIMIT ?, ? is handed LIMIT as well. Each
// engine decides what a reference becomes and how the references are
// numbered.
func Rewrite(sql string, bind func(name, lastWord string) string) string {
	var (
		out      strings.Builder
		lastWord string
		i        = 0
	)
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
			out.WriteString(bind("", lastWord))
			lastWord = afterReference(lastWord)
			i++
		case c == 's' && namedArgRe.MatchString(sql[i:]):
			m := namedArgRe.FindStringSubmatch(sql[i:])
			out.WriteString(bind(m[2], lastWord))
			lastWord = afterReference(lastWord)
			i += len(m[0])
		case c == '{' && typedParamRe.MatchString(sql[i:]):
			m := typedParamRe.FindStringSubmatch(sql[i:])
			out.WriteString(bind(m[1], lastWord))
			lastWord = afterReference(lastWord)
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
	return out.String()
}

// afterReference is the word the reference after one is preceded by: a
// LIMIT's, so that both counts of LIMIT ?, ? are bound as counts, and
// otherwise none.
func afterReference(lastWord string) string {
	if strings.EqualFold(lastWord, "limit") {
		return lastWord
	}
	return ""
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
