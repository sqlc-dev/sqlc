package sqlite

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/kyleconroy/sqlc/internal/engine/sqlite/parser"
	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
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
	blob, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	input := antlr.NewInputStream(string(blob))
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
			out := convert(stmt)
			if _, ok := out.(*ast.TODO); ok {
				continue
			}
			stmts = append(stmts, ast.Statement{
				Raw: &ast.RawStmt{
					Stmt:         out,
					StmtLocation: loc,
					StmtLen:      stmt.GetStop().GetStop() - loc + 1,
				},
			})
			loc = stmt.GetStop().GetStop() - loc
		}
	}
	return stmts, nil
}

func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntax{
		Dash: true,
	}
}
