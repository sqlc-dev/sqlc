package sqlite

import (
	"errors"
	"io"

	meyer "github.com/sqlc-dev/meyer/ast"
	"github.com/sqlc-dev/meyer/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
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
	src := string(blob)
	parsed, err := parser.ParseString(src)
	if err != nil {
		return nil, normalizeErr(err)
	}

	var stmts []ast.Statement
	// loc tracks the first byte after the previous statement's terminator, so
	// that a statement's extent covers the comments written above it. sqlc
	// reads the "-- name:" annotation out of that range.
	loc := 0
	for _, raw := range parsed {
		converter := &cc{}
		out := converter.convert(raw)
		if _, ok := out.(*ast.TODO); !ok {
			stmts = append(stmts, ast.Statement{
				Raw: &ast.RawStmt{
					Stmt:         out,
					StmtLocation: loc,
					StmtLen:      trimTerminator(src, raw) - loc,
				},
			})
		}
		loc = raw.End()
	}
	return stmts, nil
}

// trimTerminator returns the end of stmt with its terminating semicolon, and
// any space before it, removed. A statement's span runs through the
// semicolon, but sqlc's statement text does not include it.
func trimTerminator(src string, stmt meyer.Stmt) int {
	end := stmt.End()
	if end > stmt.Pos() && end <= len(src) && src[end-1] == ';' {
		end--
	}
	for end > stmt.Pos() && isSpace(src[end-1]) {
		end--
	}
	return end
}

func isSpace(b byte) bool {
	switch b {
	case ' ', '\t', '\n', '\r', '\f':
		return true
	}
	return false
}

// normalizeErr turns a meyer parse failure into the positioned error sqlc
// reports to the user. meyer's message and byte offset match SQLite's own.
func normalizeErr(err error) error {
	var perr *parser.Error
	if !errors.As(err, &perr) {
		return err
	}
	line, column := parser.LineCol(perr.SQL, perr.Offset)
	if line == 0 {
		return errors.New(perr.Message)
	}
	return &sqlerr.Error{
		Message: perr.Message,
		Line:    line,
		Column:  column,
	}
}

func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		Hash:      false,
		SlashStar: true,
	}
}
