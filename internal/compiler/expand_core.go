package compiler

import (
	"strings"

	"github.com/sqlc-dev/sqlc/internal/core"
	"github.com/sqlc-dev/sqlc/internal/source"
	"github.com/sqlc-dev/sqlc/internal/sql/ast"
)

// expandCore rewrites the stars in a query's text with the columns the core
// analyzer resolved them to. The analyzer has already walked the statement and
// its subqueries, so there is nothing left to look up here: what remains is
// deciding how each name is written, which is the engine's business and not
// the core's.
func (c *Compiler) expandCore(raw *ast.RawStmt, stars []core.StarExpansion) []source.Edit {
	if len(stars) == 0 {
		return nil
	}
	edits := make([]source.Edit, 0, len(stars))
	seen := make(map[int]bool, len(stars))
	for _, star := range stars {
		// The analyzer types an expression again when a later clause refers to
		// the output name it was given, so "SELECT (SELECT * FROM t) AS x ...
		// GROUP BY x" reports the star in the subquery once per pass. The
		// passes agree on what the star covers, and editing the same reference
		// twice would overlap, so only the first is kept.
		if seen[star.Location] {
			continue
		}
		seen[star.Location] = true

		// Everything before the star qualifies it: "foo.*" is scoped to foo,
		// while a bare "*" covers every relation in the FROM clause.
		scope := strings.Join(star.Fields[:len(star.Fields)-1], ".")

		// An unqualified star that covers more than one relation may name the
		// same column twice, so those are written with their relation.
		counts := map[string]int{}
		if scope == "" {
			for _, col := range star.Columns {
				counts[col.Name]++
			}
		}

		cols := make([]string, 0, len(star.Columns))
		for _, col := range star.Columns {
			cname := col.Name
			if star.Alias != "" {
				cname = star.Alias
			}
			cname = c.quoteIdent(cname)
			if scope != "" {
				cname = c.quoteIdent(scope) + "." + cname
			}
			if counts[cname] > 1 {
				cname = c.quoteIdent(col.Relation) + "." + cname
			}

			// This is important for SQLite in particular which needs to wrap
			// jsonb column values with `json(colname)` so they're in a publicly
			// usable format (i.e. not jsonb).
			cols = append(cols, c.selector.ColumnExpr(cname, &Column{
				Name:     col.Name,
				DataType: col.DataType,
			}))
		}

		edits = append(edits, source.Edit{
			Location: star.Location - raw.StmtLocation,
			OldFunc:  c.starOldFunc(star.Fields),
			New:      strings.Join(cols, ", "),
		})
	}
	return edits
}
