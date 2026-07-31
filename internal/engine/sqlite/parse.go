package sqlite

import (
	"errors"
	"io"
	"strings"

	meyer "github.com/sqlc-dev/meyer/ast"
	"github.com/sqlc-dev/meyer/lexer"
	"github.com/sqlc-dev/meyer/parser"
	"github.com/sqlc-dev/meyer/token"

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
	src, folded := foldSqlcCalls(string(blob))
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
		converter := &cc{folded: folded}
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

// sqlcSchema is the pseudo-schema sqlc's own functions are written under.
const sqlcSchema = "sqlc"

// foldSqlcCalls rewrites "sqlc.name(" to "sqlc_name(" and reports the offset
// of every identifier it folded.
//
// SQLite has no schema-qualified function call, so sqlc.arg(), sqlc.narg(),
// sqlc.slice() and sqlc.embed() are not SQL that meyer will parse. Folding
// the dot into the name substitutes one byte for another, which leaves every
// offset in the tree pointing at the same place in the original source;
// convertFuncCall restores the schema for the calls listed here. Working from
// the token stream rather than the raw text keeps the rewrite out of string
// literals and comments.
func foldSqlcCalls(src string) (string, map[int]bool) {
	toks := lexer.Lex(src)
	var out []byte
	var folded map[int]bool
	for i := 0; i+3 < len(toks); i++ {
		if toks[i].Kind != token.ID || toks[i+1].Kind != token.DOT || toks[i+3].Kind != token.LP {
			continue
		}
		if name := toks[i+2].Kind; name != token.ID && !token.CanFallback(name) {
			continue
		}
		// The three tokens have to be written as one word for the fold to
		// produce one. Anything else is left to fail as the syntax error it
		// is, rather than silently reparsed as something else.
		if toks[i].End != toks[i+1].Pos || toks[i+1].End != toks[i+2].Pos {
			continue
		}
		if !strings.EqualFold(src[toks[i].Pos:toks[i].End], sqlcSchema) {
			continue
		}
		if out == nil {
			out = []byte(src)
			folded = map[int]bool{}
		}
		out[toks[i+1].Pos] = '_'
		folded[toks[i].Pos] = true
	}
	if out == nil {
		return src, nil
	}
	return string(out), folded
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
