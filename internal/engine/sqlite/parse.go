package sqlite

import (
	"errors"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite/parser"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type errorListener struct {
	*antlr.DefaultErrorListener

	err string
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol interface{}, line, column int, msg string, e antlr.RecognitionException) {
	el.err = msg
}

// func (el *errorListener) ReportAmbiguity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, exact bool, ambigAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
// }
//
// func (el *errorListener) ReportAttemptingFullContext(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex int, conflictingAlts *antlr.BitSet, configs antlr.ATNConfigSet) {
// }
//
// func (el *errorListener) ReportContextSensitivity(recognizer antlr.Parser, dfa *antlr.DFA, startIndex, stopIndex, prediction int, configs antlr.ATNConfigSet) {
// }

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
	input := antlr.NewInputStream(src)
	lexer := parser.NewSQLiteLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	pp := parser.NewSQLiteParser(stream)
	el := &errorListener{}
	pp.AddErrorListener(el)
	// pp.BuildParseTrees = true
	tree := pp.Parse()
	if el.err != "" {
		return nil, errors.New(el.err)
	}
	pctx, ok := tree.(*parser.ParseContext)
	if !ok {
		return nil, fmt.Errorf("expected ParserContext; got %T\n", tree)
	}

	// ANTLR's InputStream operates on characters (runes), so token
	// positions are character indices.  source.Pluck slices with byte
	// offsets.  Build a lookup table so we can translate correctly when
	// the input contains multi-byte UTF-8 characters (e.g. em-dash).
	runeToByteOffset := buildRuneToByteOffsets(src)

	var stmts []ast.Statement
	for _, istmt := range pctx.AllSql_stmt_list() {
		list, ok := istmt.(*parser.Sql_stmt_listContext)
		if !ok {
			return nil, fmt.Errorf("expected Sql_stmt_listContext; got %T\n", istmt)
		}
		loc := 0

		for _, stmt := range list.AllSql_stmt() {
			converter := &cc{}
			out := converter.convert(stmt)
			if _, ok := out.(*ast.TODO); ok {
				loc = stmt.GetStop().GetStop() + 2
				continue
			}
			byteLoc := runeToByteOffset[loc]
			byteEnd := runeToByteOffset[stmt.GetStop().GetStop()+1]
			stmts = append(stmts, ast.Statement{
				Raw: &ast.RawStmt{
					Stmt:         out,
					StmtLocation: byteLoc,
					StmtLen:      byteEnd - byteLoc,
				},
			})
			loc = stmt.GetStop().GetStop() + 2
		}
	}
	return stmts, nil
}

// buildRuneToByteOffsets returns a slice mapping rune index to byte offset.
// Entry i holds the byte offset where rune i begins; the final entry holds
// len(s) so that an exclusive end position can be looked up safely.
func buildRuneToByteOffsets(s string) []int {
	n := utf8.RuneCountInString(s)
	offsets := make([]int, 0, n+1)
	for bytePos := range s {
		offsets = append(offsets, bytePos)
	}
	offsets = append(offsets, len(s))
	return offsets
}

func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		Hash:      false,
		SlashStar: true,
	}
}
