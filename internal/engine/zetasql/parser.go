package zetasql

import (
	"io"

	zsql "github.com/goccy/go-zetasql"

	"github.com/kyleconroy/sqlc/internal/debug"
	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	loc := zsql.NewParseResumeLocation(string(blob))
	opt := zsql.NewParserOptions()

	var stmts []ast.Statement

	converter := &cc{}

	for {
		stmt, eof, err := zsql.ParseNextStatement(loc, opt)

		debug.Dump(stmt)

		if err != nil {
			return nil, err
		}

		out := converter.convert(stmt)
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt: out,
				//StmtLocation: loc,
				// StmtLen:      stmtLen,
			},
		})

		if eof {
			break
		}
	}

	// 	stmts = append(stmts, ast.Statement{
	// 		Raw: &ast.RawStmt{
	// 			Stmt:         out,
	// 			StmtLocation: loc,
	// 			StmtLen:      stmtLen,
	// 		},
	// 	})
	return stmts, nil
}

// https://dev.mysql.com/doc/refman/8.0/en/comments.html
func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntax{
		Dash:      true,
		SlashStar: true,
		Hash:      true,
	}
}

func (p *Parser) IsReservedKeyword(s string) bool {
	return false
}
