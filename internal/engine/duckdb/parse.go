package duckdb

import (
	"context"
	"errors"
	"io"

	"github.com/sqlc-dev/darkwing/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct{}

// Parse converts DuckDB SQL into sqlc's AST. darkwing assigns every
// statement a byte span and the spans tile the input — each statement
// starts where the previous one ended — so a statement's span covers the
// comments before it, which is where the sqlc query annotation lives.
func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	stmts, err := parser.Parse(context.Background(), r)
	if err != nil {
		var perr *parser.Error
		if errors.As(err, &perr) {
			serr := &sqlerr.Error{
				Message: err.Error(),
				Line:    perr.Line,
				Column:  perr.Column,
			}
			if perr.Offset >= 0 {
				serr.Location = perr.Offset
			}
			return nil, serr
		}
		return nil, err
	}

	var out []ast.Statement
	for _, stmt := range stmts {
		converter := &cc{}
		node := converter.convert(stmt)
		if _, ok := node.(*ast.TODO); ok {
			continue
		}
		out = append(out, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         node,
				StmtLocation: stmt.Pos(),
				StmtLen:      stmt.End() - stmt.Pos(),
			},
		})
	}
	return out, nil
}

// https://duckdb.org/docs/stable/sql/introduction#comments
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true, // -- comment
		SlashStar: true, // /* comment */
	}
}
