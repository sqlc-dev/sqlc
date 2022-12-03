package cockroach

import (
	"fmt"
	"io"
	"strings"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"
	crparser "github.com/cockroachdb/cockroachdb-parser/pkg/sql/parser"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	// ctx := context.Background()
	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	body := string(contents)
	idx := 0

	var stmts []ast.Statement
	for {
		pos, ok := parser.SplitFirstStatement(body)
		if !ok {
			break
		}

		text := body[:pos]
		node, err := crparser.ParseOne(text)
		if err != nil {
			return nil, err
		}

		n := convert(node.AST)
		if n == nil {
			return nil, fmt.Errorf("unexpected nil node")
		}

		loc := strings.Index(body, text)

		stmtLen := len(text)
		if text[stmtLen-1] == ';' {
			stmtLen -= 1 // Subtract one to remove semicolon
		}

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         n,
				StmtLocation: idx + loc,
				StmtLen:      stmtLen,
			},
		})

		body = body[pos:]
		idx += pos
	}

	return stmts, nil
}

// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-COMMENTS
func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntax{
		Dash:      true,
		SlashStar: true,
	}
}
