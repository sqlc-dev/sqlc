package googlesql

import (
	"github.com/sqlc-dev/zetajones/token"
)

// IsReservedKeyword reports whether s is a reserved keyword in GoogleSQL; see
// https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#reserved_keywords
func (p *Parser) IsReservedKeyword(s string) bool {
	return token.IsReservedKeyword(s)
}
