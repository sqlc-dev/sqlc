package sqlite

import (
	"errors"
	"fmt"
	"io"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/engine/sqlite/parser"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type errorListener struct {
	*antlr.DefaultErrorListener

	err string
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
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
	// antlr reads the input as a slice of runes and reports all token
	// positions as rune indices. The rest of sqlc treats query offsets as byte
	// offsets into the original source (see source.Pluck / source.Mutate), so
	// any multi-byte rune in the source would otherwise shift every later
	// offset and corrupt the generated SQL. Build a rune-index -> byte-offset
	// table so we can translate antlr's positions back to byte offsets.
	src := string(blob)
	runeToByte := runeOffsets(src)
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
	var stmts []ast.Statement
	for _, istmt := range pctx.AllSql_stmt_list() {
		list, ok := istmt.(*parser.Sql_stmt_listContext)
		if !ok {
			return nil, fmt.Errorf("expected Sql_stmt_listContext; got %T\n", istmt)
		}
		loc := 0

		for _, stmt := range list.AllSql_stmt() {
			converter := &cc{convertPos: runeToByte}
			out := converter.convert(stmt)
			if _, ok := out.(*ast.TODO); ok {
				loc = stmt.GetStop().GetStop() + 2
				continue
			}
			end := stmt.GetStop().GetStop() + 1
			// antlr reports rune-based offsets; translate them to byte offsets
			// before they reach the byte-oriented compiler.
			byteLoc := runeToByte(loc)
			byteLen := runeToByte(end) - byteLoc
			stmts = append(stmts, ast.Statement{
				Raw: &ast.RawStmt{
					Stmt:         out,
					StmtLocation: byteLoc,
					StmtLen:      byteLen,
				},
			})
			loc = stmt.GetStop().GetStop() + 2
		}
	}
	return stmts, nil
}

// runeOffsets returns a function that maps a rune index in src (as reported by
// antlr) to the corresponding byte offset. Indices at or past the end of the
// input map to len(src) so that end-exclusive bounds keep working.
func runeOffsets(src string) func(int) int {
	// table[i] is the byte offset of the i-th rune; the final entry is len(src).
	table := make([]int, 0, len(src)+1)
	for i := range src {
		table = append(table, i)
	}
	table = append(table, len(src))
	return func(runeIdx int) int {
		if runeIdx < 0 {
			return runeIdx
		}
		if runeIdx >= len(table) {
			return len(src)
		}
		return table[runeIdx]
	}
}

func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		Hash:      false,
		SlashStar: true,
	}
}
