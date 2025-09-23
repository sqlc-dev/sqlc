package ydb

import (
	"errors"
	"fmt"
	"io"

	"github.com/antlr4-go/antlr/v4"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	parser "github.com/ydb-platform/yql-parsers/go"
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
	content := string(blob)
	input := antlr.NewInputStream(content)
	lexer := parser.NewYQLLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	pp := parser.NewYQLParser(stream)
	el := &errorListener{}
	pp.AddErrorListener(el)
	// pp.BuildParseTrees = true
	tree := pp.Sql_query()
	if el.err != "" {
		return nil, errors.New(el.err)
	}
	pctx, ok := tree.(*parser.Sql_queryContext)
	if !ok {
		return nil, fmt.Errorf("expected ParserContext; got %T\n ", tree)
	}
	var stmts []ast.Statement
	stmtListCtx := pctx.Sql_stmt_list()
	if stmtListCtx != nil {
		loc := 0
		for _, stmt := range stmtListCtx.AllSql_stmt() {
			converter := &cc{content: string(blob)}
			out, ok := stmt.Accept(converter).(ast.Node)
			if !ok {
				return nil, fmt.Errorf("expected ast.Node; got %T", out)
			}
			if _, ok := out.(*ast.TODO); ok {
				loc = byteOffset(content, stmt.GetStop().GetStop() + 2)
				continue
			}
			if out != nil {
				len := byteOffset(content, stmt.GetStop().GetStop() + 1) - loc
				stmts = append(stmts, ast.Statement{
					Raw: &ast.RawStmt{
						Stmt:         out,
						StmtLocation: loc,
						StmtLen:      len,
					},
				})
				loc = byteOffset(content, stmt.GetStop().GetStop() + 2)
			}
		}
	}
	return stmts, nil
}

func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		Hash:      false,
		SlashStar: true,
	}
}
