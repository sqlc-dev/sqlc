package sqlite

// QuoteIdent returns a quoted identifier if it needs quoting.
// SQLite uses double quotes for quoting identifiers (SQL standard),
// though backticks are also supported for MySQL compatibility.
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
