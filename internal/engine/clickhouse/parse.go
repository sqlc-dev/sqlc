package clickhouse

import (
	"bytes"
	"context"
	"io"

	"github.com/sqlc-dev/doubleclick/parser"

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

	var stmts []ast.Statement
	for _, stmt := range stmtNodes {
		converter := &cc{}
		out := converter.convert(stmt)
		if _, ok := out.(*ast.TODO); ok {
			continue
		}

		// Get position information from the statement
		pos := stmt.Pos()
		end := stmt.End()
		stmtLen := end.Offset - pos.Offset

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: pos.Offset,
				StmtLen:      stmtLen,
			},
		})
	}

	return stmts, nil
}

// https://clickhouse.com/docs/en/sql-reference/syntax#comments
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,  // -- comment
		SlashStar: true,  // /* comment */
		Hash:      true,  // # comment (ClickHouse supports this)
	}
}
