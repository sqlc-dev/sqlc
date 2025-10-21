package duckdb

import (
	"io"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// NewParser creates a new DuckDB parser
// DuckDB uses database-backed validation, so this parser is minimal
// All actual parsing and validation happens in the database via the analyzer
func NewParser() *Parser {
	return &Parser{}
}

type Parser struct{}

// Parse returns a minimal AST for DuckDB
// Since DuckDB uses database-backed catalog and analyzer,
// we don't need to parse SQL into a detailed AST.
// The analyzer will send queries to the database for validation.
func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// Return a single TODO statement containing the raw SQL
	// The database will parse and validate this later
	return []ast.Statement{
		{
			Raw: &ast.RawStmt{
				Stmt:         &ast.TODO{},
				StmtLocation: 0,
				StmtLen:      len(blob),
			},
		},
	}, nil
}

// https://duckdb.org/docs/sql/dialect/syntax#comments
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,  // -- comments
		SlashStar: true,  // /* */ comments
		Hash:      false, // DuckDB doesn't support # comments
	}
}
