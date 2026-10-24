package duckdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/sqlc-dev/darkwing/parser"

	"github.com/sqlc-dev/sqlc/internal/metadata"
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
	src, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	stmts, err := parser.Parse(context.Background(), bytes.NewReader(src))
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
			// A statement with no query annotation (INSTALL, LOAD, SET in a
			// schema file) is skipped. A statement with one is a query that
			// the user expects sqlc to generate: dropping it would hide the
			// missing method, so report it.
			span := string(src[stmt.Pos():stmt.End()])
			name, _, merr := metadata.ParseQueryNameAndType(span, metadata.CommentSyntax(p.CommentSyntax()))
			if merr != nil {
				return nil, &sqlerr.Error{Message: merr.Error(), Location: stmt.Pos() + leadingSpace(span)}
			}
			if name != "" {
				return nil, &sqlerr.Error{
					Message:  fmt.Sprintf("query %s: statement is not supported by the duckdb engine", name),
					Location: stmt.Pos() + leadingSpace(span),
				}
			}
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

// leadingSpace returns the length of the whitespace before the first
// character of s, so an error points at the annotation and not at the blank
// lines that follow the previous statement.
func leadingSpace(s string) int {
	return len(s) - len(strings.TrimLeft(s, " \t\r\n"))
}

// https://duckdb.org/docs/stable/sql/introduction#comments
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true, // -- comment
		SlashStar: true, // /* comment */
	}
}
