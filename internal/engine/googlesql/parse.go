package googlesql

import (
	"bytes"
	"context"
	"io"

	"github.com/sqlc-dev/zetajones/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	stmtNodes, err := parser.Parse(ctx, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	// zetajones node locations are byte offsets into the input. Following the
	// convention used by the other engines, a statement's location starts at the
	// end of the previous statement (or the start of the input) so that any
	// leading comment — in particular the "-- name:" annotation sqlc relies on —
	// is captured as part of the statement.
	var stmts []ast.Statement
	loc := 0
	for _, stmt := range stmtNodes {
		converter := &cc{}
		out := converter.convert(stmt)
		if _, ok := out.(*ast.TODO); ok {
			// Skip over the unsupported statement (and its trailing semicolon)
			// so the next statement's leading comment is measured from here.
			loc = stmt.End() + 1
			continue
		}

		end := stmt.End()
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: loc,
				StmtLen:      end - loc,
			},
		})
		loc = end + 1
	}

	return stmts, nil
}

// https://cloud.google.com/bigquery/docs/reference/standard-sql/lexical#comments
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true, // -- comment
		SlashStar: true, // /* comment */
		Hash:      true, // # comment
	}
}
