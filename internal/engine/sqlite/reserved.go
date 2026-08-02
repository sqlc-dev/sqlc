package sqlite

import "github.com/sqlc-dev/meyer/token"

// IsReservedKeyword reports whether s is a keyword in SQLite; see
// https://sqlite.org/lang_keywords.html. The table lives in meyer, which
// transcribes it from SQLite's own aKeywordTable.
func (p *Parser) IsReservedKeyword(s string) bool {
	return token.Lookup(s) != token.ID
}
