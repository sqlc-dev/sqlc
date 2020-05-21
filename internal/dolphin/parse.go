package dolphin

import (
	"io"
	"io/ioutil"

	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"

	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

func NewParser() *Parser {
	return &Parser{parser.New()}
}

type Parser struct {
	pingcap *parser.Parser
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	stmtNodes, _, err := p.pingcap.Parse(string(blob), "", "")
	if err != nil {
		return nil, err
	}
	var stmts []ast.Statement
	for i := range stmtNodes {
		out := convert(stmtNodes[i])
		if _, ok := out.(*ast.TODO); ok {
			continue
		}
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{Stmt: out},
			// TODO: StmtLocation
			// TODO: StmtLen
		})
	}
	return stmts, nil
}
