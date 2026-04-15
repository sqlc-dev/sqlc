package compiler

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/config"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
	"github.com/sqlc-dev/sqlc/internal/sql/astutils"
)

func (c *Compiler) expand(qc *QueryCatalog, raw *ast.RawStmt) ([]source.Edit, error) {
	// Return early if there are no A_Star nodes to expand
	stars := astutils.Search(raw, func(node ast.Node) bool {
		_, ok := node.(*ast.A_Star)
		return ok
	})
	if len(stars.Items) == 0 {
		return nil, nil
	}
	list := astutils.Search(raw, func(node ast.Node) bool {
		switch node.(type) {
		case *ast.DeleteStmt:
		case *ast.InsertStmt:
		case *ast.SelectStmt:
		case *ast.UpdateStmt:
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

var validPostgresIdent = regexp.MustCompile(`^[a-z_][a-z0-9_$]*$`)

func (c *Compiler) quoteIdent(ident string) string {
	if c.parser.IsReservedKeyword(ident) {
		return c.quote(ident)
	}
	// SQL identifiers and key words must begin with a letter (a-z, but also
	// letters with diacritical marks and non-Latin letters) or an underscore
	// (_). Subsequent characters in an identifier or key word can be letters,
	// underscores, digits (0-9), or dollar signs ($).
	//
	// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-IDENTIFIERS
	if c.conf.Engine == config.EnginePostgreSQL {
		// camelCase means the column is also camelCase
		if strings.ToLower(ident) != ident {
			return c.quote(ident)
		}
		if !validPostgresIdent.MatchString(strings.ToLower(ident)) {
			return c.quote(ident)
		}
	}
	return ident
}

func (c *Compiler) quote(x string) string {
	switch c.conf.Engine {
	case config.EngineMySQL:
		return "`" + x + "`"
	default:
		return "\"" + x + "\""
	}
}

func (c *Compiler) expandStmt(qc *QueryCatalog, raw *ast.RawStmt, node ast.Node) ([]source.Edit, error) {
	tables, err := c.sourceTables(qc, node)
	if err != nil {
		return nil, err
	}

	var targets *ast.List
	switch n := node.(type) {
	case *ast.DeleteStmt:
		targets = n.ReturningList
	case *ast.InsertStmt:
		targets = n.ReturningList
	case *ast.SelectStmt:
		targets = n.TargetList
	case *ast.UpdateStmt:
		targets = n.ReturningList
	default:
		return nil, fmt.Errorf("outputColumns: unsupported node type: %T", n)
	}

	var edits []source.Edit
	for _, target := range targets.Items {
		res, ok := target.(*ast.ResTarget)
		if !ok {
			continue
		}
		ref, ok := res.Val.(*ast.ColumnRef)
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
			case *ast.A_Star:
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

				// This is important for SQLite in particular which needs to
				// wrap jsonb column values with `json(colname)` so they're in a
				// publicly usable format (i.e. not jsonb).
				cname = c.selector.ColumnExpr(cname, column)
				cols = append(cols, cname)
			}
		}
		var old []string
		for _, p := range parts {
			if p == "*" {
				old = append(old, p)
			} else {
				old = append(old, c.quoteIdent(p))
			}
		}

		var oldString string
		var oldFunc func(string) int

		// use the sqlc.embed string instead
		if embed, ok := qc.embeds.Find(ref); ok {
			oldString = embed.Orig()
		} else {
			oldFunc = func(s string) int {
				length := 0
				for i, o := range old {
					if hasSeparator := i > 0; hasSeparator {
						length++
					}
					if strings.HasPrefix(s[length:], o) {
						length += len(o)
					} else if quoted := c.quote(o); strings.HasPrefix(s[length:], quoted) {
						length += len(quoted)
					} else {
						length += len(o)
					}
				}
				return length
			}
		}

		edits = append(edits, source.Edit{
			Location: res.Location - raw.StmtLocation,
			Old:      oldString,
			OldFunc:  oldFunc,
			New:      strings.Join(cols, ", "),
		})
	}

	return edits, nil
}
