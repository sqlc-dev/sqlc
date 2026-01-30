package clickhouse

// QuoteIdent returns a quoted identifier if it needs quoting.
// ClickHouse uses backticks or double quotes for quoting identifiers.
func (p *Parser) QuoteIdent(s string) string {
	// For now, don't quote - can be extended to quote when necessary
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
// ClickHouse uses {name:Type} for named parameters, but for positional
// parameters we use ? which is supported by the clickhouse-go driver.
func (p *Parser) Param(n int) string {
	return "?"
}

// NamedParam returns the named parameter placeholder for the given name.
// ClickHouse uses {name:Type} syntax for named parameters.
func (p *Parser) NamedParam(name string) string {
	return "{" + name + ":String}"
}

// Cast returns a type cast expression.
// ClickHouse uses CAST(expr AS type) syntax, same as MySQL.
func (p *Parser) Cast(arg, typeName string) string {
	return "CAST(" + arg + " AS " + typeName + ")"
}
