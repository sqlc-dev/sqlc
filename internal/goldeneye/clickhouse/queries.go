package clickhouse

import (
	"fmt"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

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

// bindPlaceholders rewrites sqlc's parameter syntax (?, sqlc.arg(name),
// sqlc.narg(name)) into constants ClickHouse can analyze. ClickHouse binds
// every ? positionally, so each placeholder is its own parameter even when
// a name repeats, which is how sqlc numbers them too.
func bindPlaceholders(sql string) (string, []placeholder) {
	var phs []placeholder
	out := endtoend.Rewrite(sql, func(name, lastWord string) string {
		phs = append(phs, placeholder{Number: len(phs) + 1, Name: name})
		return sentinelFor(lastWord, len(phs))
	})
	return out, phs
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
