package compiler

import (
	"fmt"
	"strings"

	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/source"
	"github.com/kyleconroy/sqlc/internal/sql/ast"
	"github.com/kyleconroy/sqlc/internal/sql/ast/pg"
	"github.com/kyleconroy/sqlc/internal/sql/astutils"
	"github.com/kyleconroy/sqlc/internal/sql/lang"
)

func (c *Compiler) expand(qc *QueryCatalog, raw *ast.RawStmt) ([]source.Edit, error) {
	list := astutils.Search(raw, func(node ast.Node) bool {
		switch node.(type) {
		case *pg.DeleteStmt:
		case *pg.InsertStmt:
		case *pg.SelectStmt:
		case *pg.UpdateStmt:
		default:
			return false
		}
		return true
	})
	if len(list.Items) == 0 {
		return nil, nil
	}
	var edits []source.Edit
	for _, item := range list.Items {
		edit, err := c.expandStmt(qc, raw, item)
		if err != nil {
			return nil, err
		}
		edits = append(edits, edit...)
	}
	return edits, nil
}

func (c *Compiler) quoteIdent(ident string) string {
	// TODO: Add a method to the parser / engine for this instead
	if lang.IsReservedKeyword(ident) {
		switch c.conf.Engine {
		case config.EngineMySQL, config.EngineMySQLBeta:
			return "`" + ident + "`"
		default:
			return "\"" + ident + "\""
		}
	}
	return ident
}

func (c *Compiler) expandStmt(qc *QueryCatalog, raw *ast.RawStmt, node ast.Node) ([]source.Edit, error) {
	tables, err := sourceTables(qc, node)
	if err != nil {
		return nil, err
	}

	var targets *ast.List
	switch n := node.(type) {
	case *pg.DeleteStmt:
		targets = n.ReturningList
	case *pg.InsertStmt:
		targets = n.ReturningList
	case *pg.SelectStmt:
		targets = n.TargetList
	case *pg.UpdateStmt:
		targets = n.ReturningList
	default:
		return nil, fmt.Errorf("outputColumns: unsupported node type: %T", n)
	}

	var edits []source.Edit
	for _, target := range targets.Items {
		res, ok := target.(*pg.ResTarget)
		if !ok {
			continue
		}
		ref, ok := res.Val.(*pg.ColumnRef)
		if !ok {
			continue
		}
		if !hasStarRef(ref) {
			continue
		}
		var parts, cols []string
		for _, f := range ref.Fields.Items {
			switch field := f.(type) {
			case *ast.String:
				parts = append(parts, field.Str)
			case *pg.String:
				parts = append(parts, field.Str)
			case *pg.A_Star:
				parts = append(parts, "*")
			default:
				return nil, fmt.Errorf("unknown field in ColumnRef: %T", f)
			}
		}
		scope := astutils.Join(ref.Fields, ".")
		counts := map[string]int{}
		if scope == "" {
			for _, t := range tables {
				for _, c := range t.Columns {
					counts[c.Name] += 1
				}
			}
		}
		for _, t := range tables {
			if scope != "" && scope != t.Rel.Name {
				continue
			}
			tableName := c.quoteIdent(t.Rel.Name)
			scopeName := c.quoteIdent(scope)
			for _, column := range t.Columns {
				cname := column.Name
				if res.Name != nil {
					cname = *res.Name
				}
				cname = c.quoteIdent(cname)
				if scope != "" {
					cname = scopeName + "." + cname
				}
				if counts[cname] > 1 {
					cname = tableName + "." + cname
				}
				cols = append(cols, cname)
			}
		}
		var old []string
		for _, p := range parts {
			old = append(old, c.quoteIdent(p))
		}
		edits = append(edits, source.Edit{
			Location: res.Location - raw.StmtLocation,
			Old:      strings.Join(old, "."),
			New:      strings.Join(cols, ", "),
		})
	}
	return edits, nil
}
