package runtime

import (
	"errors"

	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/catalog"
	"github.com/kyleconroy/sqlc/internal/sql/validate"
)

type Query struct {
}

var ErrUnsupportedStatementType = errors.New("parseQuery: unsupported statement type")

func parseQuery(c *catalog.Catalog, stmt ast.Node, src string, rewriteParameters bool) (*Query, error) {
	if err := validate.ParamStyle(stmt); err != nil {
		return nil, err
	}
	if err := validate.ParamRef(stmt); err != nil {
		return nil, err
	}
	raw, ok := stmt.(*ast.RawStmt)
	if !ok {
		return nil, errors.New("node is not a statement")
	}
	switch n := raw.Stmt.(type) {
	case *pg.SelectStmt:
	case *pg.DeleteStmt:
	case *pg.InsertStmt:
		if err := validate.InsertStmt(n); err != nil {
			return nil, err
		}
	case *pg.TruncateStmt:
	case *pg.UpdateStmt:
	default:
		return nil, ErrUnsupportedStatementType
	}

	rawSQL, err := source.Pluck(src, raw.StmtLocation, raw.StmtLen)
	if err != nil {
		return nil, err
	}
	if rawSQL == "" {
		return nil, errors.New("missing semicolon at end of file")
	}

	return &Query{}, nil
}
