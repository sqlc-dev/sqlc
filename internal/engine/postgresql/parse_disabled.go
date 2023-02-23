//go:build windows || !cgo
// +build windows !cgo

package postgresql

import (
	"errors"
	"io"
	"runtime"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	if runtime.GOOS == "windows" {
		return nil, errors.New("the PostgreSQL engine does not support Windows.")
	}
	return nil, errors.New("the PostgreSQL engine requires cgo. Please set CGO_ENABLED=1.")
}

// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-COMMENTS
func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntax{
		Dash:      true,
		SlashStar: true,
	}
}
