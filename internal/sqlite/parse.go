package sqlite

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/antlr/antlr4/runtime/Go/antlr"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sqlite/parser"
)

type listener struct {
	*parser.BaseSQLiteListener

	stmt *ast.RawStmt

	stmts []ast.Statement
}

func (l *listener) EnterSql_stmt(c *parser.Sql_stmtContext) {
	l.stmt = nil
}

func (l *listener) ExitSql_stmt(c *parser.Sql_stmtContext) {
	if l.stmt != nil {
		l.stmts = append(l.stmts, ast.Statement{
			Raw: l.stmt,
		})
		return
	}
}

func (l *listener) EnterCreate_table_stmt(c *parser.Create_table_stmtContext) {
	name := ast.TableName{
		Name: c.Table_name().GetText(),
	}

	if c.Database_name() != nil {
		name.Schema = c.Database_name().GetText()
	}

	stmt := &ast.CreateTableStmt{
		Name:        &name,
		IfNotExists: c.K_EXISTS() != nil,
	}

	for _, idef := range c.AllColumn_def() {
		if def, ok := idef.(*parser.Column_defContext); ok {
			stmt.Cols = append(stmt.Cols, &ast.ColumnDef{
				Colname: def.Column_name().GetText(),
				TypeName: &ast.TypeName{
					Name: def.Type_name().GetText(),
				},
			})
		}
	}

	l.stmt = &ast.RawStmt{Stmt: stmt}
}

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
	tree := pp.Parse()
	if el.err != "" {
		return nil, errors.New(el.err)
	}
	// p.BuildParseTrees = true
	antlr.ParseTreeWalkerDefault.Walk(l, tree)
	return l.stmts, nil
}
