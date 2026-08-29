package dolphin

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/sqlc-dev/marino/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/sqlerr"
)

func NewParser() *Parser {
	return &Parser{pingcap: parser.New()}
}

// NewFormatParser returns the parser sqlc fmt uses. It differs from the
// compiler's parser in one way: identifiers keep the case the author wrote.
// The compiler lowercases them so catalog lookups are case-insensitive, but
// a formatter that prints `Event` as `event` renames a table on the many
// servers where table names are case-sensitive.
func NewFormatParser() *Parser {
	return &Parser{pingcap: parser.New(), preserveCase: true}
}

type Parser struct {
	pingcap      *parser.Parser
	preserveCase bool
}

var lineColumn = regexp.MustCompile(`^line (\d+) column (\d+) (.*)`)

func normalizeErr(err error) error {
	if err == nil {
		return err
	}
	parts := strings.Split(err.Error(), "\n")
	msg := strings.TrimSpace(parts[0] + "\"")
	out := lineColumn.FindStringSubmatch(msg)
	if len(out) == 4 {
		line, lineErr := strconv.Atoi(out[1])
		col, colErr := strconv.Atoi(out[2])
		if lineErr != nil || colErr != nil {
			return errors.New(msg)
		}
		return &sqlerr.Error{
			Message: "syntax error",
			Err:     errors.New(out[3]),
			Line:    line,
			Column:  col,
		}
	}
	return errors.New(msg)
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	f, err := p.ParseFile(r)
	if err != nil {
		return nil, err
	}
	// The compiler skips statements sqlc has no node for; the formatter
	// must see them, so the filter lives here, not in ParseFile.
	var stmts []ast.Statement
	for _, stmt := range f.Stmts {
		if _, ok := stmt.Raw.Stmt.(*ast.TODO); ok {
			continue
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// ParseFile parses like Parse and also carries the file's comments, which
// marino's lexer records as it scans.
func (p *Parser) ParseFile(r io.Reader) (*ast.File, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	src := string(blob)
	stmtNodes, _, err := p.pingcap.Parse(src, "", "")
	if err != nil {
		return nil, normalizeErr(err)
	}
	var stmts []ast.Statement
	// A statement's text spans from the end of the previous statement
	// through its terminator, so it carries the comments written above it
	// (that's where the "-- name:" annotation lives). Each text is a
	// contiguous slice of src laid down after the one before it, so
	// searching from the previous statement's end pins every text to its
	// own occurrence even when two statements read the same.
	searchFrom := 0
	for i := range stmtNodes {
		converter := &cc{preserveCase: p.preserveCase}
		// A statement sqlc has no node for converts to a TODO and stays in
		// the list: the formatter needs its extent to keep it as written,
		// and Parse filters it out for the compiler.
		out := converter.convert(stmtNodes[i])

		text := stmtNodes[i].Text()
		idx := strings.Index(src[searchFrom:], text)
		if idx < 0 {
			return nil, fmt.Errorf("could not locate statement %d in source", i)
		}
		loc := searchFrom + idx
		searchFrom = loc + len(text)

		stmtLen := len(text)
		if stmtLen > 0 && text[stmtLen-1] == ';' {
			stmtLen -= 1 // Subtract one to remove semicolon
		}

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: loc,
				StmtLen:      stmtLen,
			},
		})
	}

	var comments []ast.Comment
	for _, c := range p.pingcap.Comments() {
		comments = append(comments, ast.Comment{
			Text:    strings.TrimRight(src[c.Begin:c.End], " \t\r\n"),
			Start:   c.Begin,
			End:     c.End,
			OwnLine: ownLine(src, c.Begin),
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

// https://dev.mysql.com/doc/refman/8.0/en/comments.html
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		SlashStar: true,
		Hash:      true,
	}
}
