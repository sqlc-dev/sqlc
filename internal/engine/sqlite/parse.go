package sqlite

import (
	"errors"
	"io"
	"strings"

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

// parseOptions accepts ORDER BY and LIMIT on UPDATE and DELETE. SQLite has
// them only when built with SQLITE_ENABLE_UPDATE_DELETE_LIMIT, which meyer
// cannot tell from the SQL, so the caller has to ask. sqlc has always
// accepted them.
var parseOptions = parser.Options{UpdateDeleteLimit: true}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	src := string(blob)
	parsed, err := parseOptions.ParseString(terminateNamedQueries(src))
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

// terminateNamedQueries adds a virtual terminator before a subsequent sqlc
// query annotation when the preceding query omitted one. sqlc annotations
// delimit queries for the other engines, and a query file may therefore
// contain multiple valid queries even though the complete SQLite script would
// otherwise be invalid. Replacing the newline immediately before the
// annotation preserves every byte offset reported by the parser.
func terminateNamedQueries(src string) string {
	terminated := []byte(src)
	for lineStart := 0; lineStart < len(terminated); {
		lineEnd := strings.IndexByte(src[lineStart:], '\n')
		if lineEnd < 0 {
			lineEnd = len(terminated)
		} else {
			lineEnd += lineStart
		}
		if lineStart > 0 && strings.HasPrefix(src[lineStart:lineEnd], "-- name: ") {
			i := lineStart - 1
			for i >= 0 && isSpace(terminated[i]) {
				i--
			}
			if i >= 0 && terminated[i] != ';' {
				terminated[lineStart-1] = ';'
			}
		}
		lineStart = lineEnd + 1
	}
	return string(terminated)
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
