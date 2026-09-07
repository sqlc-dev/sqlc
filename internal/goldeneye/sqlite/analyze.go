package sqlite

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// The analyze cases are checked against the analysis shell, which is asked
// three things about each query. What the library reports about the
// prepared statement, through `.stats stmt`: each result column's name and,
// for one read straight from a table, its declared type and which table
// column it is, with the column's NOT NULL looked up in the catalog. What
// the query returns, over the case's fixture and again over no rows, for
// the columns the library has no answer for: an expression's type is the
// storage class of its value, since SQLite types values rather than
// expressions, and any column is nullable when a NULL comes back — an
// aggregate over no rows, the far side of an outer join. And the bytecode
// the statement compiles to, which is where the parameters are found,
// since the library reports nothing about them but their number.

// placeholder is one parameter of a query as sqlc numbers them: each ? in
// turn, and each sqlc.arg name once, at its first appearance. SQLite
// numbers ?NNN the same way, so the query is rewritten with those.
type placeholder struct {
	Number int
	Name   string
}

func bind(sql string) (string, []placeholder) {
	var phs []placeholder
	numbers := map[string]int{}
	out := endtoend.Rewrite(sql, func(name, _ string) string {
		n, ok := numbers[name]
		if name == "" || !ok {
			n = len(phs) + 1
			phs = append(phs, placeholder{Number: n, Name: name})
			if name != "" {
				numbers[name] = n
			}
		}
		return fmt.Sprintf("?%d", n)
	})
	return out, phs
}

// Analyze runs a case's queries through the analysis shell and returns what
// SQLite reports in the JSON shape sqlc analyze prints.
func Analyze(ctx context.Context, dir string, c endtoend.Case) ([]byte, error) {
	if err := checkOptions(ctx, dir, analysis); err != nil {
		return nil, err
	}
	binary := analysis.binary(dir)
	schema, err := os.ReadFile(c.Schema)
	if err != nil {
		return nil, err
	}
	var fixture []byte
	if c.Fixture != "" {
		if fixture, err = os.ReadFile(c.Fixture); err != nil {
			return nil, err
		}
	}
	queries, err := c.Queries()
	if err != nil {
		return nil, err
	}
	out := make([]endtoend.AnalyzedQuery, 0, len(queries))
	for _, q := range queries {
		aq, err := analyzeQuery(ctx, binary, string(schema), string(fixture), q)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", q.Name, err)
		}
		out = append(out, aq)
	}
	return endtoend.Encode(out)
}

// Check compares what SQLite reports for a case with the output the case
// committed, returning a diff when they differ.
func Check(ctx context.Context, dir string, c endtoend.Case) (string, error) {
	got, err := Analyze(ctx, dir, c)
	if err != nil {
		return "", err
	}
	return c.Compare(got)
}

func analyzeQuery(ctx context.Context, binary, schema, fixture string, q endtoend.Query) (endtoend.AnalyzedQuery, error) {
	sql, phs := bind(q.SQL)

	// What the library says: the catalog, the bytecode, and the statement
	// prepared, with nothing bound to its parameters.
	s := newScript()
	s.sql(schema)
	s.sql(fixture)
	s.add(".mode json")
	s.add(".explain off")
	for _, cq := range catalogQueries {
		s.section(cq.section)
		s.sql(cq.sql)
	}
	s.section("explain")
	s.sql("EXPLAIN " + sql)
	s.section("query")
	s.add(".stats stmt")
	line := s.sql(sql)
	s.section("end")
	out, err := run(ctx, binary, s)
	if err != nil {
		return endtoend.AnalyzedQuery{}, err
	}
	if errs := out.errorsBefore(line); len(errs) > 0 {
		return endtoend.AnalyzedQuery{}, errors.New(strings.Join(errs, "\n"))
	}
	query := out.sections["query"]
	if query == nil || !query.prepared {
		msg := strings.TrimSpace(out.stderr)
		if msg == "" {
			msg = "the statement was not prepared"
		}
		return endtoend.AnalyzedQuery{}, errors.New(msg)
	}
	cat, err := readCatalog(out)
	if err != nil {
		return endtoend.AnalyzedQuery{}, err
	}
	var prog []instr
	if explain := out.sections["explain"]; explain != nil && len(explain.blocks) > 0 {
		if err := explain.decode(0, &prog); err != nil {
			return endtoend.AnalyzedQuery{}, fmt.Errorf("reading the bytecode: %w", err)
		}
	}
	var names []string
	for _, m := range query.columns {
		names = append(names, m.Name)
	}
	t := newTracer(cat, names)
	t.run(prog)
	params := make([]endtoend.AnalyzedColumn, len(phs))
	for i, ph := range phs {
		params[i] = t.param(ph.Number)
	}

	// What the statement returns: over the fixture, with each parameter
	// bound to a value of the column it stands in for, so that the rows
	// come through the WHERE clause, and over no rows.
	s = newScript()
	s.sql(schema)
	s.sql(fixture)
	s.add(".parameter init")
	for i, ph := range phs {
		if v := sample(params[i]); v != "" {
			s.sql(fmt.Sprintf("REPLACE INTO temp.sqlite_parameters(key, value) VALUES ('?%d', %s)", ph.Number, v))
		}
	}
	s.add(".mode quote")
	s.section("query")
	s.sql(sql)
	s.section("end")
	bound, err := run(ctx, binary, s)
	if err != nil {
		return endtoend.AnalyzedQuery{}, err
	}
	s = newScript()
	s.sql(schema)
	s.add(".mode quote")
	s.section("query")
	s.sql(sql)
	s.section("end")
	empty, err := run(ctx, binary, s)
	if err != nil {
		return endtoend.AnalyzedQuery{}, err
	}
	var rows []string
	for _, o := range []*output{bound, empty} {
		if q := o.sections["query"]; q != nil {
			rows = append(rows, q.rows...)
		}
	}

	aq := endtoend.AnalyzedQuery{
		Name:    q.Name,
		Cmd:     q.Cmd,
		Columns: []endtoend.AnalyzedColumn{},
		Params:  []endtoend.AnalyzedParam{},
	}
	for i, m := range query.columns {
		aq.Columns = append(aq.Columns, describeColumn(cat, m, classes(rows, i)))
	}
	for i, ph := range phs {
		ac := params[i]
		if ph.Name != "" {
			ac.Name = ph.Name
		}
		aq.Params = append(aq.Params, endtoend.AnalyzedParam{Number: ph.Number, Column: ac})
	}
	return aq, nil
}

// sample is an expression for a value to bind to a parameter: one of its
// column's values in the fixture, or a value of its type when all that is
// known is the type. Empty when nothing is known.
func sample(ac endtoend.AnalyzedColumn) string {
	if ac.Table != "" && ac.Name != "" {
		col, tbl := quoteIdent(ac.Name), quoteIdent(ac.Table)
		return fmt.Sprintf("(SELECT %s FROM %s WHERE %s IS NOT NULL LIMIT 1)", col, tbl, col)
	}
	if ac.Type == nil {
		return ""
	}
	switch ac.Type.Name {
	case "integer", "numeric":
		return "1"
	case "real":
		return "1.0"
	case "text":
		return "'a'"
	case "blob":
		return "x'00'"
	}
	return ""
}

func quoteIdent(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
}

// classes lists the storage classes of a column's values across rows.
func classes(rows []string, i int) []string {
	var out []string
	for _, row := range rows {
		if c := cells(row); i < len(c) {
			out = append(out, storageClass(c[i]))
		}
	}
	return out
}

// describeColumn describes a result column: as the table column it is read
// from when it is one, otherwise by the storage class of its values, and
// nullable when any of its values was NULL.
func describeColumn(cat *catalog, m columnMeta, classes []string) endtoend.AnalyzedColumn {
	ac := endtoend.AnalyzedColumn{Name: m.Name}
	if col := cat.lookup(m.Table, m.Origin); col != nil {
		d := col.describe()
		ac.Type, ac.Table = d.Type, d.Table
	} else {
		for _, class := range classes {
			if class != "null" {
				ac.Type = &endtoend.TypeExpr{Name: class}
				break
			}
		}
	}
	if ac.Type != nil {
		for _, class := range classes {
			if class == "null" {
				ac.Type.Nullable = true
			}
		}
	}
	return ac
}
