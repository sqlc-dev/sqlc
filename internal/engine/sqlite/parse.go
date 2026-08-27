package sqlite

import (
	"errors"
	"io"
	"strings"

	meyer "github.com/sqlc-dev/meyer/ast"
	"github.com/sqlc-dev/meyer/parser"
	mtoken "github.com/sqlc-dev/meyer/token"

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
	f, err := p.ParseFile(r)
	if err != nil {
		return nil, err
	}
	// The compiler skips statements sqlc has no node for (PRAGMA and
	// friends); the formatter must see them, so the filter lives here,
	// not in ParseFile.
	var stmts []ast.Statement
	for _, stmt := range f.Stmts {
		if _, ok := stmt.Raw.Stmt.(*ast.TODO); ok {
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// ParseFile parses like Parse and also carries the file's comments, taken
// from the trivia meyer's single lexer pass produces alongside the tokens.
func (p *Parser) ParseFile(r io.Reader) (*ast.File, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	src := string(blob)
	parsed, err := parseOptions.ParseFile(src)
	if err != nil {
		return nil, normalizeErr(err)
	}

	var stmts []ast.Statement
	// loc tracks the first byte after the previous statement's terminator, so
	// that a statement's extent covers the comments written above it. sqlc
	// reads the "-- name:" annotation out of that range.
	loc := 0
	for _, raw := range parsed.Stmts {
		converter := &cc{src: src}
		// A statement sqlc has no node for converts to a TODO and stays in
		// the list: the formatter needs its extent to keep it as written,
		// and Parse filters it out for the compiler.
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         converter.convert(raw),
				StmtLocation: loc,
				StmtLen:      trimTerminator(src, raw) - loc,
			},
		})
		loc = raw.End()
	}

	var comments []ast.Comment
	for _, tr := range parsed.Trivia {
		if tr.Kind != mtoken.COMMENT {
			continue
		}
		comments = append(comments, ast.Comment{
			Text:    strings.TrimRight(tr.Text(src), " \t\r\n"),
			Start:   tr.Pos,
			End:     tr.End,
			OwnLine: ownLine(src, tr.Pos),
		})
	}
	return &ast.File{Stmts: stmts, Comments: comments}, nil
}

// ownLine reports that only blank space sits between the preceding line
// break and pos.
func ownLine(src string, pos int) bool {
	for j := pos - 1; j >= 0; j-- {
		switch src[j] {
		case '\n':
			return true
		case ' ', '\t', '\r':
			continue
		default:
			return false
		}
	}
	return true
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
