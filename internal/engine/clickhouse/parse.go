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
	loc := 0
	for _, stmt := range stmtNodes {
		start := stmt.Pos().Offset - 1
		if start < loc {
			start = loc
		}
		end := statementEnd(blob, start)

		converter := &cc{}
		out := converter.convert(stmt)
		if _, ok := out.(*ast.TODO); ok {
			loc = end
			continue
		}

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: loc,
				StmtLen:      end - loc,
			},
		})
		loc = end
	}

	return stmts, nil
}

func statementEnd(blob []byte, start int) int {
	for i := start; i < len(blob); i++ {
		switch blob[i] {
		case '\'', '"', '`':
			i = skipQuoted(blob, i)
		case '-':
			if i+1 < len(blob) && blob[i+1] == '-' {
				i = skipLineComment(blob, i)
			}
		case '#':
			i = skipLineComment(blob, i)
		case '/':
			if i+1 < len(blob) && blob[i+1] == '*' {
				i = skipBlockComment(blob, i)
			}
		case ';':
			return i + 1
		}
	}
	return len(blob)
}

func skipQuoted(blob []byte, i int) int {
	q := blob[i]
	for j := i + 1; j < len(blob); j++ {
		switch blob[j] {
		case '\\':
			j++
		case q:
			if j+1 < len(blob) && blob[j+1] == q {
				j++
				continue
			}
			return j
		}
	}
	return len(blob) - 1
}

func skipLineComment(blob []byte, i int) int {
	for j := i; j < len(blob); j++ {
		if blob[j] == '\n' {
			return j
		}
	}
	return len(blob) - 1
}

func skipBlockComment(blob []byte, i int) int {
	for j := i + 2; j < len(blob); j++ {
		if blob[j] == '*' && j+1 < len(blob) && blob[j+1] == '/' {
			return j + 1
		}
	}
	return len(blob) - 1
}

// https://clickhouse.com/docs/en/sql-reference/syntax#comments
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true, // -- comment
		SlashStar: true, // /* comment */
		Hash:      true, // # comment (ClickHouse supports this)
	}
}
