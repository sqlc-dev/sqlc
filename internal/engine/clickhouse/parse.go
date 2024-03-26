package clickhouse

import (
	"io"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (c *Parser) Parse(io.Reader) ([]ast.Statement, error) {

	return []ast.Statement{}, nil
}
func (c *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{}

}
func (c *Parser) IsReservedKeyword(string) bool {
	return false
}
