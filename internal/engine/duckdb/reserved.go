package duckdb

import "github.com/sqlc-dev/darkwing/token"

// IsReservedKeyword reports whether s is a reserved keyword in DuckDB's
// grammar, per the keyword lists darkwing vendors from DuckDB.
func (p *Parser) IsReservedKeyword(s string) bool {
	return token.Categories(s).Has(token.KeywordReserved)
}
