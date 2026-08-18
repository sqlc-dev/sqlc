package mssql

import (
	"bytes"
	"context"
	"io"
	"unicode/utf8"

	"github.com/sqlc-dev/teesql/parser"

	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{}
}

type Parser struct{}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	script, err := parser.Parse(ctx, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}

	toByte := byteOffsets(blob)

	var stmts []ast.Statement
	loc := 0
	for _, batch := range script.Batches {
		for _, stmt := range batch.Statements {
			start := loc
			if frag, ok := stmt.(fragmented); ok && frag.Frag().HasSpan() {
				if s := toByte(frag.Frag().StartOffset); s > start {
					start = s
				}
			}
			end := statementEnd(blob, start)

			converter := &cc{}
			out := converter.convert(stmt)
			if _, ok := out.(*ast.TODO); ok {
				loc = end
				continue
			}

			stmts = append(stmts, ast.Statement{
				Raw: &ast.RawStmt{
					Stmt:         out,
					StmtLocation: loc,
					StmtLen:      end - loc,
				},
			})
			loc = end
		}
	}

	return stmts, nil
}

// byteOffsets returns a function mapping a UTF-16 code-unit offset — how
// teesql records source spans, mirroring ScriptDom — to a byte offset in
// blob. For ASCII input the two are identical and the identity is returned.
func byteOffsets(blob []byte) func(int) int {
	ascii := true
	for _, b := range blob {
		if b >= utf8.RuneSelf {
			ascii = false
			break
		}
	}
	if ascii {
		return func(off int) int { return off }
	}
	// Walk the blob once, recording the byte index at which each UTF-16
	// offset begins.
	byteAt := make([]int, 0, len(blob)+1)
	for i := 0; i < len(blob); {
		r, size := utf8.DecodeRune(blob[i:])
		units := 1
		if r > 0xFFFF {
			units = 2
		}
		for u := 0; u < units; u++ {
			byteAt = append(byteAt, i)
		}
		i += size
	}
	byteAt = append(byteAt, len(blob))
	return func(off int) int {
		if off < 0 {
			return 0
		}
		if off >= len(byteAt) {
			return len(blob)
		}
		return byteAt[off]
	}
}

// statementEnd scans from start for the semicolon ending the statement,
// skipping string literals, quoted and bracketed identifiers, and comments.
func statementEnd(blob []byte, start int) int {
	for i := start; i < len(blob); i++ {
		switch blob[i] {
		case '\'', '"':
			i = skipQuoted(blob, i)
		case '[':
			i = skipBracketed(blob, i)
		case '-':
			if i+1 < len(blob) && blob[i+1] == '-' {
				i = skipLineComment(blob, i)
			}
		case '/':
			if i+1 < len(blob) && blob[i+1] == '*' {
				i = skipBlockComment(blob, i)
			}
		case ';':
			return i + 1
		}
	}
	return len(blob)
}

// skipQuoted skips a string literal or quoted identifier; T-SQL escapes the
// quote character by doubling it.
func skipQuoted(blob []byte, i int) int {
	q := blob[i]
	for j := i + 1; j < len(blob); j++ {
		if blob[j] == q {
			if j+1 < len(blob) && blob[j+1] == q {
				j++
				continue
			}
			return j
		}
	}
	return len(blob) - 1
}

// skipBracketed skips a [bracketed identifier]; a closing bracket is escaped
// by doubling it.
func skipBracketed(blob []byte, i int) int {
	for j := i + 1; j < len(blob); j++ {
		if blob[j] == ']' {
			if j+1 < len(blob) && blob[j+1] == ']' {
				j++
				continue
			}
			return j
		}
	}
	return len(blob) - 1
}

func skipLineComment(blob []byte, i int) int {
	for j := i; j < len(blob); j++ {
		if blob[j] == '\n' {
			return j
		}
	}
	return len(blob) - 1
}

func skipBlockComment(blob []byte, i int) int {
	for j := i + 2; j < len(blob); j++ {
		if blob[j] == '*' && j+1 < len(blob) && blob[j+1] == '/' {
			return j + 1
		}
	}
	return len(blob) - 1
}

// https://learn.microsoft.com/en-us/sql/t-sql/language-elements/comment-transact-sql
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true, // -- comment
		SlashStar: true, // /* comment */
	}
}
