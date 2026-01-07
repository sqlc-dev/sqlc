package dolphin

// QuoteIdent returns a quoted identifier if it needs quoting.
// MySQL uses backticks for quoting identifiers.
func (p *Parser) QuoteIdent(s string) string {
	// For now, don't quote - MySQL is less strict about quoting
	return s
}

// TypeName returns the SQL type name for the given namespace and name.
// Handles MySQL-specific type name mappings for formatting.
func (p *Parser) TypeName(ns, name string) string {
	if ns != "" {
		return ns + "." + name
	}
	// Map internal type names to MySQL CAST-compatible names for formatting
	switch name {
	case "bigint unsigned":
		return "UNSIGNED"
	case "bigint signed":
		return "SIGNED"
	}
	return name
}

// Param returns the parameter placeholder for the given number.
// MySQL uses ? for all parameters (positional).
func (p *Parser) Param(n int) string {
	return "?"
}

// NamedParam returns the named parameter placeholder for the given name.
// MySQL doesn't have native named parameters, so we use ? (positional).
// The actual parameter names are handled by sqlc's rewrite phase.
func (p *Parser) NamedParam(name string) string {
	return "?"
}

// Cast returns a type cast expression.
// MySQL uses CAST(expr AS type) syntax.
func (p *Parser) Cast(arg, typeName string) string {
	return "CAST(" + arg + " AS " + typeName + ")"
}
