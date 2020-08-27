package dolphin

import (
	"errors"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/pingcap/parser"
	_ "github.com/pingcap/parser/test_driver"

	"github.com/kyleconroy/sqlc/internal/metadata"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/sqlerr"
)

func NewParser() *Parser {
	return &Parser{parser.New()}
}

type Parser struct {
	pingcap *parser.Parser
}

var lineColumn = regexp.MustCompile(`^line (\d+) column (\d+) (.*)`)

func normalizeErr(err error) error {
	if err == nil {
		return err
	}
	parts := strings.Split(err.Error(), "\n")
	msg := strings.TrimSpace(parts[0] + "\"")
	out := lineColumn.FindStringSubmatch(msg)
	if len(out) == 4 {
		line, lineErr := strconv.Atoi(out[1])
		col, colErr := strconv.Atoi(out[2])
		if lineErr != nil || colErr != nil {
			return errors.New(msg)
		}
		return &sqlerr.Error{
			Message: "syntax error",
			Err:     errors.New(out[3]),
			Line:    line,
			Column:  col,
		}
	}
	return errors.New(msg)
}

func (p *Parser) Parse(r io.Reader) ([]ast.Statement, error) {
	blob, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	stmtNodes, _, err := p.pingcap.Parse(string(blob), "", "")
	if err != nil {
		return nil, normalizeErr(err)
	}
	var stmts []ast.Statement
	for i := range stmtNodes {
		converter := &cc{}
		out := converter.convert(stmtNodes[i])
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
				StmtLen:      len(text) - 1, // Subtract one to remove semicolon
			},
		})
	}
	return stmts, nil
}

// https://dev.mysql.com/doc/refman/8.0/en/comments.html
func (p *Parser) CommentSyntax() metadata.CommentSyntax {
	return metadata.CommentSyntax{
		Dash:      true,
		SlashStar: true,
		Hash:      true,
	}
}
