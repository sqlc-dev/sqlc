package clickhouse

// QuoteIdent returns a quoted identifier if it needs quoting.
// ClickHouse uses backticks for quoting identifiers.
func (p *Parser) QuoteIdent(s string) string {
	// For now, don't quote - return as-is
	return s
}

// TypeName returns the SQL type name for the given namespace and name.
func (p *Parser) TypeName(ns, name string) string {
	if ns != "" {
		return ns + "." + name
	}
	return name
}

// Param returns the parameter placeholder for the given number.
// ClickHouse uses ? for positional parameters.
func (p *Parser) Param(n int) string {
	return "?"
}

// NamedParam returns the named parameter placeholder for the given name.
// ClickHouse uses @name syntax for named parameters.
func (p *Parser) NamedParam(name string) string {
	return "@" + name
}

// Cast returns a type cast expression.
// ClickHouse uses CAST(expr AS type) syntax.
func (p *Parser) Cast(arg, typeName string) string {
	return "CAST(" + arg + " AS " + typeName + ")"
}
