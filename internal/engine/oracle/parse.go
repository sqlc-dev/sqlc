package oracle

import (
	"errors"
	"fmt"
	"io"

	"github.com/antlr4-go/antlr/v4"

	"github.com/sqlc-dev/sqlc/internal/engine/oracle/parser"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

type errorListener struct {
	*antlr.DefaultErrorListener

	err string
}

func (el *errorListener) SyntaxError(recognizer antlr.Recognizer, offendingSymbol any, line, column int, msg string, e antlr.RecognitionException) {
	if el.err == "" {
		el.err = fmt.Sprintf("line %d:%d: %s", line, column, msg)
	}
}

// NewParser returns a parser for Oracle SQL / PL-SQL.
func NewParser() *Parser {
	return &Parser{}
}

type Parser struct {
}

// Parse reads Oracle SQL from r and returns the sqlc shared-AST statements it
// was able to convert. Statements the converter does not yet understand are
// represented as *ast.TODO and skipped, matching the SQLite engine's behaviour
// so the engine can ship incrementally.
func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	input := antlr.NewInputStream(string(blob))
	lexer := parser.NewPlSqlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	pp := parser.NewPlSqlParser(stream)
	el := &errorListener{}
	pp.RemoveErrorListeners()
	pp.AddErrorListener(el)

	tree := pp.Sql_script()
	if el.err != "" {
		return nil, errors.New(el.err)
	}

	script, ok := tree.(*parser.Sql_scriptContext)
	if !ok {
		return nil, fmt.Errorf("expected *parser.Sql_scriptContext; got %T", tree)
	}

	var stmts []ast.Statement
	loc := 0
	for _, iunit := range script.AllUnit_statement() {
		unit, ok := iunit.(*parser.Unit_statementContext)
		if !ok {
			continue
		}
		converter := &cc{}
		out := converter.convert(unit)
		if _, ok := out.(*ast.TODO); ok {
			loc = unit.GetStop().GetStop() + 2
			continue
		}
		length := (unit.GetStop().GetStop() + 1) - loc
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: loc,
				StmtLen:      length,
			},
		})
		loc = unit.GetStop().GetStop() + 2
	}
	return stmts, nil
}

// CommentSyntax reports the comment styles Oracle SQL supports: -- line comments
// and /* */ block comments. Oracle does not use # comments.
func (p *Parser) CommentSyntax() source.CommentSyntax {
	return source.CommentSyntax{
		Dash:      true,
		Hash:      false,
		SlashStar: true,
	}
}
