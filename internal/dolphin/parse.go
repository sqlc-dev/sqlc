package dolphin

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pingcap/parser"
	_ "github.com/pingcap/tidb/types/parser_driver"

	"github.com/kyleconroy/sqlc/internal/metadata"
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
		spew.Dump(err)
		return nil, err
	}
	var stmts []ast.Statement
	for i := range stmtNodes {
		out := convert(stmtNodes[i])
		if _, ok := out.(*ast.TODO); ok {
			continue
		}

		// TODO: Attach the text directly to the ast.Statement node
		text := stmtNodes[i].Text()
		loc := strings.Index(string(blob), text)

		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: loc,
				StmtLen:      len(text),
			},
		})
	}
	return stmts, nil
}

func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntaxStar
}
