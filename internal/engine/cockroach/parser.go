package cockroach

import (
	"fmt"
	"io"

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
	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	nodes, err := crparser.Parse(string(contents))
	if err != nil {
		return nil, err
	}
	var stmts []ast.Statement
	for _, raw := range nodes {
		n := convert(raw.AST)
		if n == nil {
			return nil, fmt.Errorf("unexpected nil node")
		}
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         n,
				StmtLocation: 0,            // TODO
				StmtLen:      len(raw.SQL), // TODO
			},
		})
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
