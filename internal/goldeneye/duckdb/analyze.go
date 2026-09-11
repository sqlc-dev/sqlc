package duckdb

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// The analyze cases are checked against the DuckDB CLI, which runs each
// case's schema and fixture into an in-memory database of their own, one
// process per question, and is asked four things about each query. What
// its parameters are: the query is prepared and explained with a string
// sentinel bound to each parameter, and the unoptimized logical plan the
// CLI prints first shows each as CAST('goldeneye_k' AS T), T being the
// type the binder gave the parameter. What its result columns are:
// DESCRIBE, with each parameter replaced by a NULL of that type, names and
// types them. Which column each is read from and which column a parameter
// stands in for: DuckDB prints a plan with every column by its bare name,
// so these are read from the query text, the select list's items and the
// operand beside each parameter resolved against the FROM clause and the
// catalog, from which a column read from a table takes its declared type
// and nullability. And whether an expression can be NULL, which DuckDB
// does not track: the query is run, with each parameter bound to a value
// of its type, over the fixture and over no rows, and a column is nullable
// when either run returns a NULL for it.

// placeholder is one parameter of a query as sqlc numbers them: a $n by
// its number, and each $name, sqlc.arg name or ? in turn at its first
// appearance, taking the lowest number no $n took.
type placeholder struct {
	Number int
	Name   string
}

var (
	numberedRe = regexp.MustCompile(`^\$([0-9]+)`)
	namedRe    = regexp.MustCompile(`^\$([A-Za-z_][A-Za-z0-9_]*)`)
	sqlcArgRe  = regexp.MustCompile(`^sqlc\.(n?arg|slice)\(\s*'?([A-Za-z_][A-Za-z0-9_]*)'?\s*\)`)
)

// bind rewrites the query so that every parameter is a $k numbered as
// sqlc numbers it, and lists the parameters in that order.
func bind(query string) (string, []placeholder) {
	type occurrence struct {
		start, end int
		number     int    // for $n
		name       string // for $name or sqlc.arg, "" for ?
	}
	var occ []occurrence
	i := 0
	for i < len(query) {
		c := query[i]
		switch {
		case c == '\'' || c == '"':
			i = quotedEnd(query, i)
		case strings.HasPrefix(query[i:], "--"):
			end := strings.IndexByte(query[i:], '\n')
			if end < 0 {
				i = len(query)
			} else {
				i += end
			}
		case strings.HasPrefix(query[i:], "/*"):
			end := strings.Index(query[i:], "*/")
			if end < 0 {
				i = len(query)
			} else {
				i += end + 2
			}
		case c == '?':
			occ = append(occ, occurrence{start: i, end: i + 1})
			i++
		case c == '$' && numberedRe.MatchString(query[i:]):
			m := numberedRe.FindStringSubmatch(query[i:])
			n, _ := strconv.Atoi(m[1])
			occ = append(occ, occurrence{start: i, end: i + len(m[0]), number: n})
			i += len(m[0])
		case c == '$' && namedRe.MatchString(query[i:]):
			m := namedRe.FindStringSubmatch(query[i:])
			occ = append(occ, occurrence{start: i, end: i + len(m[0]), name: m[1]})
			i += len(m[0])
		case c == 's' && sqlcArgRe.MatchString(query[i:]):
			m := sqlcArgRe.FindStringSubmatch(query[i:])
			occ = append(occ, occurrence{start: i, end: i + len(m[0]), name: m[2]})
			i += len(m[0])
		default:
			i++
		}
	}
	taken := map[int]bool{}
	for _, o := range occ {
		if o.number > 0 {
			taken[o.number] = true
		}
	}
	byName := map[string]int{}
	next := 1
	assign := func(o occurrence) int {
		if o.number > 0 {
			return o.number
		}
		if n, ok := byName[o.name]; ok && o.name != "" {
			return n
		}
		for taken[next] {
			next++
		}
		n := next
		taken[n] = true
		if o.name != "" {
			byName[o.name] = n
		}
		return n
	}
	var phs []placeholder
	numbered := map[int]bool{}
	var out strings.Builder
	last := 0
	for _, o := range occ {
		n := assign(o)
		if !numbered[n] {
			numbered[n] = true
			phs = append(phs, placeholder{Number: n, Name: o.name})
		}
		out.WriteString(query[last:o.start])
		out.WriteString("$" + strconv.Itoa(n))
		last = o.end
	}
	out.WriteString(query[last:])
	for i := 1; i < len(phs); i++ {
		for j := i; j > 0 && phs[j-1].Number > phs[j].Number; j-- {
			phs[j-1], phs[j] = phs[j], phs[j-1]
		}
	}
	return out.String(), phs
}

// nullMarker is what a NULL is printed as when the CLI prints CSV, unlike
// any value a case holds.
const nullMarker = "<goldeneye:null>"

// analyzer holds the CLI a case runs through, the case's schema and
// fixture, and the catalog of the schema.
type analyzer struct {
	binary  string
	schema  string // the schema's statements, each ending in ;
	fixture string // the fixture's, or ""
	tables  map[string][]column
	enums   map[string][]string // the labels of each enum type the schema created
	// canonical names the type the dialect reports a spelling as, for
	// each alias types.jsonl lists: DuckDB spells a JSON column json,
	// which its own catalog lists as a spelling of varchar.
	canonical map[string]string
}

// readAliases reads the aliases the generated types.jsonl gives each
// type, keyed by alias.
func readAliases() (map[string]string, error) {
	dir, err := dialect.Dir(Engine)
	if err != nil {
		return nil, err
	}
	types, err := dialect.ReadTypes(dir)
	if err != nil {
		return nil, err
	}
	canonical := map[string]string{}
	for _, t := range types {
		for _, alias := range t.Aliases {
			canonical[strings.ToLower(alias)] = strings.ToLower(t.Name)
		}
	}
	return canonical, nil
}

// run executes a script in a fresh in-memory database loaded with the
// schema and, unless told otherwise, the fixture, and returns what the
// CLI printed in the mode asked for: "json" prints one JSON array per
// statement that returns rows, "csv" prints rows with a header and NULL
// as nullMarker, and "" prints the CLI's own text, which is how EXPLAIN
// draws a plan. A statement that fails ends the script with DuckDB's own
// error message.
func (a *analyzer) run(ctx context.Context, script, mode string, withFixture bool) (string, error) {
	args := []string{"-bail"}
	switch mode {
	case "json":
		args = append(args, "-json")
	case "csv":
		args = append(args, "-csv", "-header", "-nullvalue", nullMarker)
	}
	args = append(args, ":memory:")
	cmd := exec.CommandContext(ctx, a.binary, args...)
	prelude := a.schema
	if withFixture {
		prelude += a.fixture
	}
	cmd.Stdin = strings.NewReader(prelude + script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", errors.New(msg)
	}
	return stdout.String(), nil
}

// query runs one statement and decodes the rows it printed.
func (a *analyzer) query(ctx context.Context, statement string) ([]map[string]json.RawMessage, error) {
	out, err := a.run(ctx, statement+";\n", "json", true)
	if err != nil {
		return nil, err
	}
	var rows []map[string]json.RawMessage
	dec := json.NewDecoder(strings.NewReader(out))
	for {
		var rs []map[string]json.RawMessage
		err := dec.Decode(&rs)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decoding duckdb output: %w", err)
		}
		rows = append(rows, rs...)
	}
	return rows, nil
}

func str(v json.RawMessage) string {
	var s string
	json.Unmarshal(v, &s)
	return s
}

// column is what the catalog says about a column of a table.
type column struct {
	name     string
	spelling string // the type as DuckDB spells it
	typ      *analysis.TypeExpr
	nullable bool
}

// statements joins a script's statements, each ending in a semicolon.
func statements(src string) string {
	s := strings.TrimRight(strings.TrimSpace(src), ";")
	if s == "" {
		return ""
	}
	return s + ";\n"
}

// Analyze runs a case's queries through the CLI and returns what DuckDB
// reports in the JSON shape sqlc analyze prints.
func Analyze(ctx context.Context, binary string, c endtoend.Case) ([]byte, error) {
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
	a := &analyzer{binary: binary, schema: statements(string(schema)), fixture: statements(string(fixture))}
	if a.canonical, err = readAliases(); err != nil {
		return nil, err
	}
	if err := a.readCatalog(ctx); err != nil {
		return nil, err
	}
	out := make([]analysis.Query, 0, len(queries))
	for _, q := range queries {
		aq, err := a.analyzeQuery(ctx, q)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", q.Name, err)
		}
		out = append(out, aq)
	}
	return analysis.Encode(out)
}

// Check compares what DuckDB reports for a case with the output the case
// committed, returning a diff when they differ.
func Check(ctx context.Context, binary string, c endtoend.Case) (string, error) {
	got, err := Analyze(ctx, binary, c)
	if err != nil {
		return "", err
	}
	return c.Compare(got)
}

// readCatalog reads the tables and enum types the schema created.
func (a *analyzer) readCatalog(ctx context.Context) error {
	rows, err := a.query(ctx, `SELECT type_name, labels FROM duckdb_types() WHERE database_name = 'memory' AND NOT internal AND logical_type = 'ENUM' ORDER BY type_name`)
	if err != nil {
		return err
	}
	a.enums = map[string][]string{}
	for _, row := range rows {
		var labels []string
		json.Unmarshal(row["labels"], &labels)
		a.enums[strings.ToLower(str(row["type_name"]))] = labels
	}
	rows, err = a.query(ctx, `SELECT table_name, column_name, data_type, is_nullable FROM duckdb_columns() WHERE database_name = 'memory' ORDER BY table_oid, column_index`)
	if err != nil {
		return err
	}
	a.tables = map[string][]column{}
	for _, row := range rows {
		table := strings.ToLower(str(row["table_name"]))
		var nullable bool
		json.Unmarshal(row["is_nullable"], &nullable)
		spelling := str(row["data_type"])
		a.tables[table] = append(a.tables[table], column{
			name:     str(row["column_name"]),
			spelling: spelling,
			typ:      a.parseType(spelling),
			nullable: nullable,
		})
	}
	return nil
}

// lookup finds a column of a table, which DuckDB matches in any case.
func (a *analyzer) lookup(table, name string) (column, bool) {
	for _, col := range a.tables[strings.ToLower(table)] {
		if strings.EqualFold(col.name, name) {
			return col, true
		}
	}
	return column{}, false
}

// resolve finds the table a reference names in a scope: the table its
// qualifier aliases, or the one table of the scope that has the column.
func (a *analyzer) resolve(sc scope, r ref) (string, column, bool) {
	tables := sc.tables
	if sc.target.name != "" {
		tables = append([]tableRef{sc.target}, tables...)
	}
	if r.qualifier != "" {
		if sc.ctes[r.qualifier] {
			return "", column{}, false
		}
		for _, t := range tables {
			if t.alias == r.qualifier || (t.alias == "" && t.name == r.qualifier) {
				col, ok := a.lookup(t.name, r.column)
				return t.name, col, ok
			}
		}
		return "", column{}, false
	}
	var found string
	var fc column
	for _, t := range tables {
		if col, ok := a.lookup(t.name, r.column); ok {
			if found != "" && found != t.name {
				return "", column{}, false
			}
			found, fc = t.name, col
		}
	}
	return found, fc, found != ""
}

// describe reads what the catalog says about a column read from a table.
func describe(table string, col column) analysis.Column {
	return analysis.Column{Name: col.name, Type: withNullable(col.typ, col.nullable), Table: table}
}

func withNullable(t *analysis.TypeExpr, nullable bool) *analysis.TypeExpr {
	if t == nil {
		return nil
	}
	out := *t
	out.Nullable = nullable
	return &out
}

// binding is what is known about one parameter: the type the binder gave
// it, spelled as DuckDB spells it, and the column it stands in for.
type binding struct {
	spelling string
	typ      *analysis.TypeExpr
	column   *analysis.Column
}

var (
	sentinelCastRe = regexp.MustCompile(`CAST\('goldeneye_([0-9]+)' AS `)
	conversionRe   = regexp.MustCompile(`Could not convert string 'goldeneye_([0-9]+)'`)
)

// explain prepares the query, explains it with a sentinel bound to each
// parameter and returns the type the binder gave each, keyed by number.
// A parameter whose sentinel the binder converts on the spot, as an
// INSERT's VALUES are, is bound to NULL instead and reported by nothing.
func (a *analyzer) explain(ctx context.Context, sql string, phs []placeholder) (map[int]string, error) {
	null := map[int]bool{}
	for attempt := 0; attempt <= len(phs); attempt++ {
		args := make([]string, len(phs))
		for i, ph := range phs {
			if null[ph.Number] {
				args[i] = "NULL"
			} else {
				args[i] = fmt.Sprintf("'goldeneye_%d'", ph.Number)
			}
		}
		script := "PREPARE goldeneye AS " + sql + ";\nSET explain_output = 'all';\nEXPLAIN EXECUTE goldeneye(" + strings.Join(args, ", ") + ");\n"
		out, err := a.run(ctx, script, "", true)
		if err != nil {
			if m := conversionRe.FindStringSubmatch(err.Error()); m != nil {
				n, _ := strconv.Atoi(m[1])
				if !null[n] {
					null[n] = true
					continue
				}
			}
			return nil, err
		}
		return sentinelTypes(out), nil
	}
	return nil, errors.New("could not bind the parameters")
}

// sentinelTypes reads the type each sentinel is cast to out of the plan
// the CLI drew, which wraps long expressions across lines inside its
// boxes.
func sentinelTypes(plan string) map[int]string {
	var b strings.Builder
	for _, line := range strings.Split(plan, "\n") {
		if strings.ContainsAny(line, "╭╮╰╯─┬┴├") {
			continue
		}
		b.WriteString(strings.Trim(line, "│ "))
		b.WriteByte(' ')
	}
	flat := b.String()
	types := map[int]string{}
	for _, m := range sentinelCastRe.FindAllStringSubmatchIndex(flat, -1) {
		n, _ := strconv.Atoi(flat[m[2]:m[3]])
		if _, ok := types[n]; ok {
			continue
		}
		// The type runs to the parenthesis closing the CAST.
		depth := 1
		i := m[1]
		for ; i < len(flat); i++ {
			if flat[i] == '(' {
				depth++
			} else if flat[i] == ')' {
				depth--
				if depth == 0 {
					break
				}
			}
		}
		types[n] = strings.Join(strings.Fields(flat[m[1]:i]), " ")
	}
	return types
}

// analyzeQuery describes one query.
func (a *analyzer) analyzeQuery(ctx context.Context, q endtoend.Query) (analysis.Query, error) {
	sql, phs := bind(q.SQL)
	aq := analysis.Query{
		Name:    q.Name,
		Cmd:     q.Cmd,
		Columns: []analysis.Column{},
		Params:  []analysis.Param{},
	}
	t := text(tokenize(sql))
	sc := t.readScope()

	// Parameters: a partner column first, then the cast the query wraps
	// the parameter in, then the type the binder gave it.
	bindings := map[int]*binding{}
	var casts []string
	var castNumbers []int
	for _, ph := range phs {
		k := strconv.Itoa(ph.Number)
		b := &binding{}
		bindings[ph.Number] = b
		if sc.kind == "insert" {
			if pos, ok := t.valuesPosition(k); ok {
				cols := t.insertColumns()
				if cols == nil {
					for _, col := range a.tables[sc.target.name] {
						cols = append(cols, col.name)
					}
				}
				if pos < len(cols) {
					if col, ok := a.lookup(sc.target.name, cols[pos]); ok {
						c := describe(sc.target.name, col)
						b.column, b.spelling, b.typ = &c, col.spelling, col.typ
						continue
					}
				}
			}
		}
		if r, ok := t.partner(k); ok {
			if table, col, ok := a.resolve(sc, r); ok {
				c := describe(table, col)
				b.column, b.spelling, b.typ = &c, col.spelling, col.typ
				continue
			}
		}
		if typ, ok := t.castOf(k); ok {
			casts = append(casts, typ)
			castNumbers = append(castNumbers, ph.Number)
		}
	}
	if len(casts) > 0 {
		// DuckDB spells the cast's type: VARCHAR(5) is VARCHAR, mood is
		// ENUM('sad', 'ok').
		var items []string
		for i, typ := range casts {
			items = append(items, fmt.Sprintf("CAST(NULL AS %s) AS p%d", typ, i))
		}
		rows, err := a.query(ctx, "DESCRIBE SELECT "+strings.Join(items, ", "))
		if err != nil {
			return analysis.Query{}, err
		}
		for i, row := range rows {
			if i < len(castNumbers) {
				b := bindings[castNumbers[i]]
				b.spelling = str(row["column_type"])
				b.typ = a.parseType(b.spelling)
			}
		}
	}
	if len(phs) > 0 {
		types, err := a.explain(ctx, sql, phs)
		if err != nil {
			return analysis.Query{}, err
		}
		for n, spelling := range types {
			if b := bindings[n]; b != nil && b.typ == nil {
				b.spelling = spelling
				b.typ = a.parseType(spelling)
			}
		}
	}

	// Result columns.
	var columns []analysis.Column
	switch sc.kind {
	case "select":
		rows, err := a.query(ctx, "DESCRIBE "+a.substitute(sql, phs, bindings, false))
		if err != nil {
			return analysis.Query{}, err
		}
		for _, row := range rows {
			columns = append(columns, analysis.Column{Name: str(row["column_name"]), Type: a.parseType(str(row["column_type"]))})
		}
		a.attribute(sc, t.selectItems(), columns)
	default:
		// A RETURNING column is the target table's column it names;
		// DuckDB describes no DML, so an expression is typed by nothing.
		for _, it := range t.returningItems() {
			switch {
			case it.star:
				for _, col := range a.tables[sc.target.name] {
					columns = append(columns, describe(sc.target.name, col))
				}
			case it.ref != nil:
				if table, col, ok := a.resolve(sc, *it.ref); ok {
					c := describe(table, col)
					c.Name = it.name
					columns = append(columns, c)
					continue
				}
				columns = append(columns, analysis.Column{Name: it.name})
			default:
				columns = append(columns, analysis.Column{Name: it.name})
			}
		}
	}
	if len(columns) > 0 {
		nullable := a.observe(ctx, sql, phs, bindings, len(columns))
		for i := range columns {
			if columns[i].Table == "" && columns[i].Type != nil && nullable[i] {
				columns[i].Type = withNullable(columns[i].Type, true)
			}
		}
	}
	aq.Columns = append(aq.Columns, columns...)

	for _, ph := range phs {
		b := bindings[ph.Number]
		ac := analysis.Column{Type: b.typ}
		if b.column != nil {
			ac = *b.column
		}
		aq.Params = append(aq.Params, analysis.Param{Number: ph.Number, Column: ac})
	}
	return aq, nil
}

// attribute says which table each result column is read from, by
// matching the select list's items, a star expanded to its table's
// columns, against the columns DuckDB described. A list that expands to
// another number of columns than DuckDB describes attributes nothing.
func (a *analyzer) attribute(sc scope, items []item, columns []analysis.Column) {
	var origins []*analysis.Column
	for _, it := range items {
		switch {
		case it.star && it.ref == nil:
			for _, t := range sc.tables {
				for _, col := range a.tables[t.name] {
					c := describe(t.name, col)
					origins = append(origins, &c)
				}
			}
		case it.star:
			table, ok := a.aliased(sc, it.ref.qualifier)
			if !ok {
				return
			}
			for _, col := range a.tables[table] {
				c := describe(table, col)
				origins = append(origins, &c)
			}
		case it.ref != nil:
			if table, col, ok := a.resolve(sc, *it.ref); ok {
				c := describe(table, col)
				origins = append(origins, &c)
			} else {
				origins = append(origins, nil)
			}
		default:
			origins = append(origins, nil)
		}
	}
	if len(origins) != len(columns) {
		return
	}
	for i, o := range origins {
		if o != nil {
			columns[i].Type = o.Type
			columns[i].Table = o.Table
		}
	}
}

// aliased finds the table a qualifier names in a scope.
func (a *analyzer) aliased(sc scope, qualifier string) (string, bool) {
	for _, t := range sc.tables {
		if t.alias == qualifier || (t.alias == "" && t.name == qualifier) {
			if _, ok := a.tables[t.name]; ok {
				return t.name, true
			}
		}
	}
	return "", false
}

// substitute replaces each parameter with a NULL of its type, or with a
// value of its type when values is set, for DuckDB to describe or run the
// query. A parameter of unknown type becomes a bare NULL.
func (a *analyzer) substitute(sql string, phs []placeholder, bindings map[int]*binding, values bool) string {
	out := sql
	for i := len(phs) - 1; i >= 0; i-- {
		ph := phs[i]
		b := bindings[ph.Number]
		repl := "NULL"
		if b != nil && b.spelling != "" {
			repl = "CAST(NULL AS " + b.spelling + ")"
			if values {
				if v, ok := a.zero(b.typ); ok {
					repl = "CAST('" + strings.ReplaceAll(v, "'", "''") + "' AS " + b.spelling + ")"
				}
			}
		}
		out = strings.ReplaceAll(out, "$"+strconv.Itoa(ph.Number), repl)
	}
	return out
}

// observe runs the query with a value bound to each parameter, over the
// fixture and over no rows, and reports which of the first n result
// columns came back NULL in either run. A run the CLI rejects observes
// nothing.
func (a *analyzer) observe(ctx context.Context, sql string, phs []placeholder, bindings map[int]*binding, n int) []bool {
	nullable := make([]bool, n)
	statement := a.substitute(sql, phs, bindings, true) + ";\n"
	runs := []bool{true}
	if a.fixture != "" {
		runs = append(runs, false)
	}
	for _, withFixture := range runs {
		out, err := a.run(ctx, statement, "csv", withFixture)
		if err != nil {
			continue
		}
		records, err := csv.NewReader(strings.NewReader(out)).ReadAll()
		if err != nil {
			continue
		}
		for _, rec := range records[min(1, len(records)):] {
			for j, v := range rec {
				if j < n && v == nullMarker {
					nullable[j] = true
				}
			}
		}
	}
	return nullable
}
