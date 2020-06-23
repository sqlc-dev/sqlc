package vitess

import (
	"io"
	"io/ioutil"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	t := sqlparser.NewStringTokenizer(string(contents))
	var stmts []ast.Statement
	var start int
	for {
		stmt, err := sqlparser.ParseNextStrictDDL(t)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		out := convert(stmt)
		if _, ok := out.(*ast.TODO); ok {
			continue
		}
		stmts = append(stmts, ast.Statement{
			Raw: &ast.RawStmt{
				Stmt:         out,
				StmtLocation: start,
				StmtLen:      t.Position - start - 1,
			},
		})
		start = t.Position
	}
	return stmts, nil
}

func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntaxStar
}
