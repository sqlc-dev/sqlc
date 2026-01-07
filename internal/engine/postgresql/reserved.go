package postgresql

import (
	"fmt"
	"strings"
)

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

// QuoteIdent returns a quoted identifier if it needs quoting.
// This implements the format.Dialect interface.
func (p *Parser) QuoteIdent(s string) string {
	if p.IsReservedKeyword(s) || hasMixedCase(s) {
		return `"` + s + `"`
	}
	return s
}

// TypeName returns the SQL type name for the given namespace and name.
// This implements the format.Dialect interface.
func (p *Parser) TypeName(ns, name string) string {
	if ns == "pg_catalog" {
		switch name {
		case "int4":
			return "integer"
		case "int8":
			return "bigint"
		case "int2":
			return "smallint"
		case "float4":
			return "real"
		case "float8":
			return "double precision"
		case "bool":
			return "boolean"
		case "bpchar":
			return "character"
		case "timestamptz":
			return "timestamp with time zone"
		case "timetz":
			return "time with time zone"
		default:
			return name
		}
	}
	if ns != "" {
		return ns + "." + name
	}
	return name
}

// Param returns the parameter placeholder for the given number.
// PostgreSQL uses $1, $2, etc.
func (p *Parser) Param(n int) string {
	return fmt.Sprintf("$%d", n)
}

// NamedParam returns the named parameter placeholder for the given name.
// PostgreSQL/sqlc uses @name syntax.
func (p *Parser) NamedParam(name string) string {
	return "@" + name
}

// Cast returns a type cast expression.
// PostgreSQL uses expr::type syntax.
func (p *Parser) Cast(arg, typeName string) string {
	return arg + "::" + typeName
}

// https://www.postgresql.org/docs/current/sql-keywords-appendix.html
func (p *Parser) IsReservedKeyword(s string) bool {
	switch strings.ToLower(s) {
	case "all":
	case "analyse":
	case "analyze":
	case "and":
	case "any":
	case "array":
	case "as":
	case "asc":
	case "asymmetric":
	case "authorization":
	case "binary":
	case "both":
	case "case":
	case "cast":
	case "check":
	case "collate":
	case "collation":
	case "column":
	case "concurrently":
	case "constraint":
	case "create":
	case "cross":
	case "current_catalog":
	case "current_date":
	case "current_role":
	case "current_schema":
	case "current_time":
	case "current_timestamp":
	case "current_user":
	case "default":
	case "deferrable":
	case "desc":
	case "distinct":
	case "do":
	case "else":
	case "end":
	case "except":
	case "false":
	case "fetch":
	case "for":
	case "foreign":
	case "freeze":
	case "from":
	case "full":
	case "grant":
	case "group":
	case "having":
	case "ilike":
	case "in":
	case "initially":
	case "inner":
	case "intersect":
	case "into":
	case "is":
	case "isnull":
	case "join":
	case "lateral":
	case "leading":
	case "left":
	case "like":
	case "limit":
	case "localtime":
	case "localtimestamp":
	case "natural":
	case "not":
	case "notnull":
	case "null":
	case "offset":
	case "on":
	case "only":
	case "or":
	case "order":
	case "outer":
	case "overlaps":
	case "placing":
	case "primary":
	case "references":
	case "returning":
	case "right":
	case "select":
	case "session_user":
	case "similar":
	case "some":
	case "symmetric":
	case "table":
	case "tablesample":
	case "then":
	case "to":
	case "trailing":
	case "true":
	case "union":
	case "unique":
	case "user":
	case "using":
	case "variadic":
	case "verbose":
	case "when":
	case "where":
	case "window":
	case "with":
	default:
		return false
	}
	return true
}
