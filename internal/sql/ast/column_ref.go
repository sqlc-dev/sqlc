package ast

import "strings"

// sqlReservedWords is a set of SQL keywords that must be quoted when used as identifiers
var sqlReservedWords = map[string]bool{
	"all": true, "analyse": true, "analyze": true, "and": true, "any": true,
	"array": true, "as": true, "asc": true, "asymmetric": true, "authorization": true,
	"between": true, "binary": true, "both": true, "case": true, "cast": true,
	"check": true, "collate": true, "collation": true, "column": true, "concurrently": true,
	"constraint": true, "create": true, "cross": true, "current_catalog": true,
	"current_date": true, "current_role": true, "current_schema": true,
	"current_time": true, "current_timestamp": true, "current_user": true,
	"default": true, "deferrable": true, "desc": true, "distinct": true, "do": true,
	"else": true, "end": true, "except": true, "false": true, "fetch": true,
	"for": true, "foreign": true, "freeze": true, "from": true, "full": true,
	"grant": true, "group": true, "having": true, "ilike": true, "in": true,
	"initially": true, "inner": true, "intersect": true, "into": true, "is": true,
	"isnull": true, "join": true, "lateral": true, "leading": true, "left": true,
	"like": true, "limit": true, "localtime": true, "localtimestamp": true,
	"natural": true, "not": true, "notnull": true, "null": true, "offset": true,
	"on": true, "only": true, "or": true, "order": true, "outer": true,
	"overlaps": true, "placing": true, "primary": true, "references": true,
	"returning": true, "right": true, "select": true, "session_user": true,
	"similar": true, "some": true, "symmetric": true, "table": true, "tablesample": true,
	"then": true, "to": true, "trailing": true, "true": true, "union": true,
	"unique": true, "user": true, "using": true, "variadic": true, "verbose": true,
	"when": true, "where": true, "window": true, "with": true,
}

// needsQuoting returns true if the identifier is a SQL reserved word
// that needs to be quoted when used as an identifier
func needsQuoting(s string) bool {
	return sqlReservedWords[strings.ToLower(s)]
}

// hasMixedCase returns true if the string has any uppercase letters
// (identifiers with mixed case need quoting in PostgreSQL)
func hasMixedCase(s string) bool {
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			return true
		}
	}
	return false
}

// quoteIdent returns a quoted identifier if it needs quoting
func quoteIdent(s string) string {
	if needsQuoting(s) || hasMixedCase(s) {
		return `"` + s + `"`
	}
	return s
}

type ColumnRef struct {
	Name string

	// From pg.ColumnRef
	Fields   *List
	Location int
}

func (n *ColumnRef) Pos() int {
	return n.Location
}

func (n *ColumnRef) Format(buf *TrackedBuffer) {
	if n == nil {
		return
	}

	if n.Fields != nil {
		var items []string
		for _, item := range n.Fields.Items {
			switch nn := item.(type) {
			case *String:
				items = append(items, quoteIdent(nn.Str))
			case *A_Star:
				items = append(items, "*")
			}
		}
		buf.WriteString(strings.Join(items, "."))
	}
}
