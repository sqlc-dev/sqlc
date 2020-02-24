package sqlite

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sqlite/parser"
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
	l := &listener{}
	el := &errorListener{}
	pp.AddErrorListener(el)
	// pp.BuildParseTrees = true
	tree := pp.Parse()
	if el.err != "" {
		return nil, errors.New(el.err)
	}
	antlr.ParseTreeWalkerDefault.Walk(l, tree)
	return l.stmts, nil
}
