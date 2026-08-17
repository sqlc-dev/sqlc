package ast

import "github.com/sqlc-dev/sqlc/internal/sql/format"

// formatReturningOptions writes the PostgreSQL 18 RETURNING WITH (...) option
// list that renames the OLD and NEW aliases available in a RETURNING clause.
func formatReturningOptions(buf *TrackedBuffer, d format.Dialect, oldAlias, newAlias string) {
	if oldAlias == "" && newAlias == "" {
		return
	}
	buf.WriteString("WITH (")
	if oldAlias != "" {
		buf.WriteString("OLD AS ")
		buf.WriteString(d.QuoteIdent(oldAlias))
	}
	if newAlias != "" {
		if oldAlias != "" {
			buf.WriteString(", ")
		}
		buf.WriteString("NEW AS ")
		buf.WriteString(d.QuoteIdent(newAlias))
	}
	buf.WriteString(") ")
}
