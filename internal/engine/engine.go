// Package engine provides the interface and registry for database engines.
// Engines are responsible for parsing SQL statements and providing database-specific
// functionality like catalog creation, keyword checking, and comment syntax.
package engine

import (
	"io"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/catalog"
)

// Parser is the interface that wraps the basic SQL parsing methods.
// All database engines must implement this interface.
type Parser interface {
	// Parse parses SQL from the given reader and returns a slice of statements.
	Parse(io.Reader) ([]ast.Statement, error)

	// CommentSyntax returns the comment syntax supported by this engine.
	CommentSyntax() source.CommentSyntax

	// IsReservedKeyword returns true if the given string is a reserved keyword.
	IsReservedKeyword(string) bool
}

// Dialect provides database-specific formatting for SQL identifiers and expressions.
// This is used when reformatting queries for output.
type Dialect interface {
	// QuoteIdent returns a quoted identifier if it needs quoting.
	QuoteIdent(string) string

	// TypeName returns the SQL type name for the given namespace and name.
	TypeName(ns, name string) string

	// Param returns the parameter placeholder for the given number.
	// E.g., PostgreSQL uses $1, MySQL uses ?, etc.
	Param(n int) string

	// NamedParam returns the named parameter placeholder for the given name.
	NamedParam(name string) string

	// Cast returns a type cast expression.
	Cast(arg, typeName string) string
}

// Selector generates output expressions for SELECT and RETURNING statements.
// Different engines may need to wrap certain column types for proper output.
type Selector interface {
	// ColumnExpr generates output to be used in a SELECT or RETURNING
	// statement based on input column name and metadata.
	ColumnExpr(name string, dataType string) string
}

// Column represents column metadata for the Selector interface.
type Column struct {
	DataType string
}

// Engine is the main interface that database engines must implement.
// It provides factory methods for creating engine-specific components.
type Engine interface {
	// Name returns the unique name of this engine (e.g., "postgresql", "mysql", "sqlite").
	Name() string

	// Parser returns a new Parser instance for this engine.
	Parser() Parser

	// Catalog returns a new Catalog instance pre-populated with built-in types and schemas.
	Catalog() *catalog.Catalog

	// Selector returns a Selector for generating column expressions.
	// Returns nil if the engine uses the default selector.
	Selector() Selector

	// Dialect returns the Dialect for this engine.
	// Returns nil if the parser implements Dialect directly.
	Dialect() Dialect
}

// EngineFactory is a function that creates a new Engine instance.
type EngineFactory func() Engine

// DefaultSelector is a selector implementation that does the simplest possible
// pass through when generating column expressions. Its use is suitable for all
// database engines not requiring additional customization.
type DefaultSelector struct{}

// ColumnExpr returns the column name unchanged.
func (s *DefaultSelector) ColumnExpr(name string, dataType string) string {
	return name
}
