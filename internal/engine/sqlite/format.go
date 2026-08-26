package sqlite

import "strings"

// QuoteIdent quotes an identifier when printing it bare would change what it
// names: the parser folds unquoted identifiers to lower case, so any name
// that is not already a plain lower-case identifier — or that collides with
// a keyword — only survives a parse round-trip inside double quotes.
func (p *Parser) QuoteIdent(s string) string {
	if s == "" || p.IsReservedKeyword(s) || !plainIdent(s) {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

// plainIdent reports whether s parses back as itself unquoted: ASCII
// lower-case letters, digits and underscores, not starting with a digit.
func plainIdent(s string) bool {
	for i, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r == '_':
		case r >= '0' && r <= '9':
			if i == 0 {
				return false
			}
		default:
			return false
		}
	}
	return true
}

// TypeName returns the SQL type name for the given namespace and name.
func (p *Parser) TypeName(ns, name string) string {
	if ns != "" {
		return ns + "." + name
	}
	return name
}

// Param returns the parameter placeholder for the given number.
// SQLite uses ? for positional parameters.
func (p *Parser) Param(n int) string {
	return "?"
}

// NamedParam returns the named parameter placeholder for the given name.
// SQLite uses :name syntax for named parameters.
func (p *Parser) NamedParam(name string) string {
	return ":" + name
}

// Cast returns a type cast expression.
// SQLite uses CAST(expr AS type) syntax.
func (p *Parser) Cast(arg, typeName string) string {
	return "CAST(" + arg + " AS " + typeName + ")"
}
