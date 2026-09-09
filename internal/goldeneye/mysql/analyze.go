package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// The analyze cases are checked against a live server, which is asked
// three things about each query. What a driver sees: the query is run,
// with every parameter a user variable set to NULL, and each result
// column's name, type and nullability are read from the result set's
// metadata, the way go-sql-driver reports them. What the resolver made of
// it: the optimizer trace prints each query block back after name
// resolution, with every column qualified, every alias kept and every
// SELECT * expanded, which says which table each result column is read
// from and what each parameter is compared with or assigned to. And for a
// statement the trace does not expand — an INSERT ... VALUES, a
// single-table UPDATE or DELETE — the note EXPLAIN leaves, which prints
// the statement the same way. Views and derived tables are kept as the
// query wrote them rather than merged, so that a column read through one
// is reported as its column. MySQL itself reports nothing about a
// parameter but its position, so a parameter is described by its partner:
// a column's type and nullability come from information_schema, a column
// of a derived table from what its block projects, and an expression's
// from running it, over the tables it reads, as a query of its own.

// placeholder is one parameter of a query as sqlc numbers them: each ? in
// turn, and each sqlc.arg name once, at its first appearance. Each becomes
// a user variable, which MySQL prints back by name, except in a LIMIT or
// OFFSET, where MySQL allows nothing but a number.
type placeholder struct {
	Number   int
	Name     string
	Limit    bool
	Sentinel string // the variable's name, or the number
}

// limitBase is added to a LIMIT parameter's number, so that the count is
// unlike any a query would write.
const limitBase = 4000000000

func bind(sql string) (string, []placeholder) {
	var phs []placeholder
	numbers := map[string]int{}
	out := endtoend.Rewrite(sql, func(name, lastWord string) string {
		n, ok := numbers[name]
		if name == "" || !ok {
			n = len(phs) + 1
			ph := placeholder{Number: n, Name: name}
			switch strings.ToLower(lastWord) {
			case "limit", "offset":
				ph.Limit = true
				ph.Sentinel = strconv.Itoa(limitBase + n)
			default:
				ph.Sentinel = "goldeneye_" + strconv.Itoa(n)
			}
			phs = append(phs, ph)
			if name != "" {
				numbers[name] = n
			}
		}
		ph := phs[n-1]
		if ph.Limit {
			return ph.Sentinel
		}
		return "@" + ph.Sentinel
	})
	return out, phs
}

// column is what information_schema says about a column.
type column struct {
	name     string
	typ      string
	nullable bool
}

// relation is a table the catalog knows, by schema and name.
type relation struct {
	schema, name string
}

// analyzer holds the session a case runs in and the catalog: the case's
// own schema, and information_schema, whose names are matched and
// reported in lower case, since MySQL matches them in any case and the
// dialect spells them that way.
type analyzer struct {
	conn    *sql.Conn
	db      string                // the case's schema
	catalog map[relation][]column // columns in order
}

// table returns the catalog's name for a table of a schema, or "" when
// the catalog has no such table.
func (a *analyzer) table(schema, name string) (relation, bool) {
	if schema == "" {
		schema = a.db
	}
	if strings.EqualFold(schema, "information_schema") {
		schema, name = "information_schema", strings.ToLower(name)
	}
	rel := relation{schema, name}
	_, ok := a.catalog[rel]
	return rel, ok
}

var dbNameRe = regexp.MustCompile(`[^A-Za-z0-9_]+`)

// Analyze loads a case's schema and fixture into a database of their own
// on the server, runs its queries there and returns what MySQL reports in
// the JSON shape sqlc analyze prints.
func Analyze(ctx context.Context, dsn string, c endtoend.Case) ([]byte, error) {
	conn, err := open(ctx, dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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

	name := "goldeneye_" + dbNameRe.ReplaceAllString(c.Name, "_")
	db := quote(name)
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS " + db,
		"CREATE DATABASE " + db,
		"USE " + db,
		"SET SESSION optimizer_trace = 'enabled=on'",
		"SET SESSION optimizer_trace_max_mem_size = 67108864",
		// A view or derived table is kept as the query wrote it rather
		// than merged into the block that reads it, so that a column read
		// from one is reported as its column, and an information_schema
		// view is not resolved away into the data dictionary tables
		// behind it.
		"SET SESSION optimizer_switch = 'derived_merge=off,derived_condition_pushdown=off'",
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return nil, err
		}
	}
	defer conn.ExecContext(context.WithoutCancel(ctx), "DROP DATABASE IF EXISTS "+db)
	if _, err := conn.ExecContext(ctx, string(schema)); err != nil {
		return nil, fmt.Errorf("loading %s: %w", c.Schema, err)
	}
	if strings.TrimSpace(string(fixture)) != "" {
		if _, err := conn.ExecContext(ctx, string(fixture)); err != nil {
			return nil, fmt.Errorf("loading %s: %w", c.Fixture, err)
		}
	}

	a := &analyzer{conn: conn, db: name}
	if a.catalog, err = readCatalog(ctx, conn, name); err != nil {
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

// Check compares what MySQL reports for a case with the output the case
// committed, returning a diff when they differ.
func Check(ctx context.Context, dsn string, c endtoend.Case) (string, error) {
	got, err := Analyze(ctx, dsn, c)
	if err != nil {
		return "", err
	}
	return c.Compare(got)
}

const catalogQuery = `
SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME, DATA_TYPE, COLUMN_TYPE, IS_NULLABLE
FROM information_schema.COLUMNS
WHERE TABLE_SCHEMA IN (?, 'information_schema')
ORDER BY TABLE_SCHEMA, TABLE_NAME, ORDINAL_POSITION`

func readCatalog(ctx context.Context, conn *sql.Conn, db string) (map[relation][]column, error) {
	rows, err := conn.QueryContext(ctx, catalogQuery, db)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	catalog := map[relation][]column{}
	for rows.Next() {
		var schema, table, name, dataType, columnType, nullable string
		if err := rows.Scan(&schema, &table, &name, &dataType, &columnType, &nullable); err != nil {
			return nil, err
		}
		if schema == "information_schema" {
			table, name = strings.ToLower(table), strings.ToLower(name)
		}
		rel := relation{schema, table}
		catalog[rel] = append(catalog[rel], column{
			name:     name,
			typ:      typeName(dataType, columnType),
			nullable: nullable == "YES",
		})
	}
	return catalog, rows.Err()
}

// statement is everything the server said about one query.
type statement struct {
	a       *analyzer
	trace   trace
	note    *text // the statement as EXPLAIN printed it, nil without one
	aliases map[string]tableRef
	results []*sql.ColumnType // of the executed query, nil for one returning no rows
	target  *relation         // the table an INSERT ... VALUES writes
	targets []string          // its columns, in row order
}

func (a *analyzer) analyzeQuery(ctx context.Context, q endtoend.Query) (analysis.Query, error) {
	sql, phs := bind(q.SQL)

	var sets []string
	for _, ph := range phs {
		if !ph.Limit {
			sets = append(sets, "@"+ph.Sentinel+" = NULL")
		}
	}
	if len(sets) > 0 {
		if _, err := a.conn.ExecContext(ctx, "SET "+strings.Join(sets, ", ")); err != nil {
			return analysis.Query{}, err
		}
	}

	// EXPLAIN resolves and optimises the statement without running it,
	// leaving its note and its trace behind. Only the traditional format
	// leaves the note, and since MySQL 26 the default format is the tree.
	if err := drain(a.conn.QueryContext(ctx, "EXPLAIN FORMAT=TRADITIONAL "+sql)); err != nil {
		return analysis.Query{}, err
	}
	note, err := readNote(ctx, a.conn)
	if err != nil {
		return analysis.Query{}, err
	}
	tr, err := readTrace(ctx, a.conn)
	if err != nil {
		return analysis.Query{}, err
	}

	s := &statement{a: a, trace: tr, aliases: map[string]tableRef{}}
	if note != "" {
		t := tokenize(note)
		s.note = &t
		t.aliases(s.aliases)
		s.target, s.targets = s.insertTargets(t)
	}
	for _, t := range tr.blocks {
		t.aliases(s.aliases)
	}

	aq := analysis.Query{
		Name:    q.Name,
		Cmd:     q.Cmd,
		Columns: []analysis.Column{},
		Params:  []analysis.Param{},
	}
	if returnsRows(sql) {
		rows, err := a.conn.QueryContext(ctx, sql)
		if err != nil {
			return analysis.Query{}, err
		}
		s.results, err = rows.ColumnTypes()
		rows.Close()
		if err != nil {
			return analysis.Query{}, err
		}
		// Names, types and nullability are the result set's. The top
		// block's SELECT list adds which table each column is read from,
		// when the two line up. MySQL sends every size of TEXT and BLOB
		// as the one wire type and tells them apart by the length beside
		// it, which the driver keeps to itself, so a column read from a
		// table, directly or through a derived table, is spelled the way
		// the table declares it.
		var items []item
		if top, ok := tr.blocks[1]; ok {
			items = top.selectItems()
		}
		for i, ct := range s.results {
			ac := analysis.Column{Name: ct.Name(), Type: typeOf(ct)}
			if len(items) == len(s.results) {
				top := tr.blocks[1]
				if it := items[i]; it.end-it.start == 1 && top.at(it.start).kind == tkIdent {
					tok := top.at(it.start)
					if rel, ok := s.resolve(tok.parts); ok {
						ac.Table = rel.name
					}
					if col, ok := s.origin(top, tok, 0); ok {
						ac.Type.Name = col.typ
					}
				}
			}
			aq.Columns = append(aq.Columns, ac)
		}
	}

	for _, ph := range phs {
		ac := analysis.Column{}
		if t, i, ok := s.locate(ph); ok {
			ac = s.describe(ctx, t, t.partner(i), 0)
		}
		if ph.Name != "" {
			ac.Name = ph.Name
		}
		aq.Params = append(aq.Params, analysis.Param{Number: ph.Number, Column: ac})
	}
	return aq, nil
}

// drain reads a result set to the end and closes it.
func drain(rows *sql.Rows, err error) error {
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
	}
	return rows.Err()
}

// readNote returns the statement EXPLAIN printed back, from the note it
// leaves among the warnings, or "" when it left none.
func readNote(ctx context.Context, conn *sql.Conn) (string, error) {
	rows, err := conn.QueryContext(ctx, "SHOW WARNINGS")
	if err != nil {
		return "", err
	}
	defer rows.Close()
	note := ""
	for rows.Next() {
		var level, message string
		var code int
		if err := rows.Scan(&level, &code, &message); err != nil {
			return "", err
		}
		if code == 1003 {
			note = message
		}
	}
	return note, rows.Err()
}

// readTrace returns the optimizer trace of the last statement.
func readTrace(ctx context.Context, conn *sql.Conn) (trace, error) {
	rows, err := conn.QueryContext(ctx, "SELECT TRACE, MISSING_BYTES_BEYOND_MAX_MEM_SIZE FROM information_schema.OPTIMIZER_TRACE")
	if err != nil {
		return trace{}, err
	}
	defer rows.Close()
	blob := ""
	for rows.Next() {
		var missing int64
		if err := rows.Scan(&blob, &missing); err != nil {
			return trace{}, err
		}
		if missing > 0 {
			return trace{}, fmt.Errorf("the optimizer trace was truncated by %d bytes", missing)
		}
	}
	if err := rows.Err(); err != nil {
		return trace{}, err
	}
	return parseTrace(blob)
}

// returnsRows reports whether a statement produces a result set.
func returnsRows(sql string) bool {
	head := strings.ToLower(strings.TrimSpace(sql))
	if strings.HasPrefix(head, "(") {
		return true
	}
	for _, kw := range []string{"select", "with", "show", "describe", "desc", "explain", "table", "values"} {
		if strings.HasPrefix(head, kw) && (len(head) == len(kw) || !isWordByte(head[len(kw)])) {
			return true
		}
	}
	return false
}

// typeOf is a result column's type as the driver reports it, spelled the
// way a column declaration does: "bigint unsigned" rather than the
// driver's "UNSIGNED BIGINT".
func typeOf(ct *sql.ColumnType) *analysis.TypeExpr {
	name := strings.ToLower(ct.DatabaseTypeName())
	if rest, ok := strings.CutPrefix(name, "unsigned "); ok {
		name = rest + " unsigned"
	}
	nullable, _ := ct.Nullable()
	return &analysis.TypeExpr{Name: name, Nullable: nullable}
}

// resolve returns the catalog table a qualified column name is read from:
// the qualifier is an alias, or a table of the case's schema, or a table
// of the schema the name spells out.
func (s *statement) resolve(parts []string) (relation, bool) {
	if len(parts) < 2 {
		return relation{}, false
	}
	qualifier := parts[len(parts)-2]
	if ref, ok := s.aliases[qualifier]; ok {
		return s.a.table(ref.schema, ref.name)
	}
	schema := ""
	if len(parts) >= 3 {
		schema = parts[len(parts)-3]
	}
	return s.a.table(schema, qualifier)
}

// locate finds the text a placeholder was printed in and where: the
// innermost query block whose expanded_query has it, or else the note.
func (s *statement) locate(ph placeholder) (text, int, bool) {
	find := func(t text) int {
		if ph.Limit {
			return t.findNumber(ph.Sentinel)
		}
		return t.findVar(ph.Sentinel)
	}
	for _, sel := range s.trace.order() {
		t := s.trace.blocks[sel]
		if i := find(t); i >= 0 {
			return t, i, true
		}
	}
	if s.note != nil {
		if i := find(*s.note); i >= 0 {
			return *s.note, i, true
		}
	}
	return text{}, -1, false
}

// describe turns a parameter's partner into the column description sqlc
// prints for the parameter.
func (s *statement) describe(ctx context.Context, t text, p partner, depth int) analysis.Column {
	switch p.kind {
	case partnerOperand:
		return s.operand(ctx, t, p.start, p.end, depth)
	case partnerProjection:
		if s.results != nil {
			if top, ok := s.trace.blocks[1]; ok && top.src == t.src && p.index < len(s.results) {
				ct := s.results[p.index]
				return analysis.Column{Name: ct.Name(), Type: typeOf(ct)}
			}
		}
		if items := t.selectItems(); p.index < len(items) {
			return analysis.Column{Name: items[p.index].alias}
		}
	case partnerInsert:
		if s.target != nil && len(s.targets) > 0 {
			return s.column(*s.target, s.targets[p.index%len(s.targets)])
		}
	case partnerLimit:
		// A LIMIT or OFFSET count is an unsigned 64-bit integer to MySQL.
		return analysis.Column{Type: &analysis.TypeExpr{Name: "bigint unsigned"}}
	}
	return analysis.Column{}
}

// operand describes an operand: a column by what the schema says of it, a
// call by its name and what it returns, anything else by its type.
func (s *statement) operand(ctx context.Context, t text, start, end, depth int) analysis.Column {
	tok := t.at(start)
	if end-start == 1 && tok.kind == tkIdent {
		return s.columnRef(ctx, t, tok, depth)
	}
	if tok.kind == tkVar {
		return analysis.Column{}
	}
	ac := analysis.Column{}
	if tok.kind == tkWord && t.at(start+1).kind == tkLParen && t.isCall(start) && !strings.HasPrefix(tok.text, "<") {
		ac.Name = tok.text
	}
	ac.Type = s.eval(ctx, t, start, end)
	return ac
}

// columnRef describes a column reference: a column of the schema, a column
// of a derived table by the expression its block projects, or a SELECT list
// alias by the expression it names.
func (s *statement) columnRef(ctx context.Context, t text, tok token, depth int) analysis.Column {
	parts := tok.parts
	if depth > 8 {
		return analysis.Column{Name: parts[len(parts)-1]}
	}
	if len(parts) == 1 {
		for _, it := range t.selectItems() {
			if it.alias == parts[0] {
				ac := s.operand(ctx, t, it.start, it.end, depth+1)
				ac.Name = parts[0]
				return ac
			}
		}
		return analysis.Column{Name: parts[0]}
	}
	qualifier, name := parts[len(parts)-2], parts[len(parts)-1]
	if rel, ok := s.resolve(parts); ok {
		return s.column(rel, name)
	}
	if sel, ok := s.trace.derived[qualifier]; ok {
		if block, ok := s.trace.blocks[sel]; ok {
			for _, it := range block.selectItems() {
				if it.alias == name {
					ac := s.operand(ctx, block, it.start, it.end, depth+1)
					ac.Name = name
					return ac
				}
			}
		}
	}
	return analysis.Column{Name: name}
}

// origin follows a column reference to the catalog column it reads, through
// any derived table or CTE whose block projects that column as it is.
func (s *statement) origin(t text, tok token, depth int) (column, bool) {
	parts := tok.parts
	if depth > 8 || len(parts) < 2 {
		return column{}, false
	}
	qualifier, name := parts[len(parts)-2], parts[len(parts)-1]
	if rel, ok := s.resolve(parts); ok {
		if rel.schema == "information_schema" {
			name = strings.ToLower(name)
		}
		for _, col := range s.a.catalog[rel] {
			if col.name == name {
				return col, true
			}
		}
		return column{}, false
	}
	sel, ok := s.trace.derived[qualifier]
	if !ok {
		return column{}, false
	}
	block, ok := s.trace.blocks[sel]
	if !ok {
		return column{}, false
	}
	for _, it := range block.selectItems() {
		if it.alias == name && it.end-it.start == 1 && block.at(it.start).kind == tkIdent {
			return s.origin(block, block.at(it.start), depth+1)
		}
	}
	return column{}, false
}

// column describes a column of a catalog table.
func (s *statement) column(rel relation, name string) analysis.Column {
	if rel.schema == "information_schema" {
		name = strings.ToLower(name)
	}
	for _, col := range s.a.catalog[rel] {
		if col.name == name {
			return analysis.Column{
				Name:  name,
				Type:  &analysis.TypeExpr{Name: col.typ, Nullable: col.nullable},
				Table: rel.name,
			}
		}
	}
	return analysis.Column{Name: name, Table: rel.name}
}

// eval finds the type of an expression by running it as a query of its
// own, over the tables its columns are read from under the aliases the
// statement gave them. An expression that reads a derived table cannot be
// run that way and gets no type.
func (s *statement) eval(ctx context.Context, t text, start, end int) *analysis.TypeExpr {
	var from []string
	seen := map[string]bool{}
	for _, id := range t.idents(start, end) {
		qualifier := id.parts[len(id.parts)-2]
		if seen[qualifier] {
			continue
		}
		seen[qualifier] = true
		rel, ok := s.resolve(id.parts)
		if !ok {
			return nil
		}
		item := quote(rel.schema) + "." + quote(rel.name)
		if qualifier != rel.name {
			item += " " + quote(qualifier)
		}
		from = append(from, item)
	}
	query := "SELECT " + t.slice(start, end)
	if len(from) > 0 {
		query += " FROM " + strings.Join(from, ", ")
	}
	rows, err := s.a.conn.QueryContext(ctx, query+" LIMIT 0")
	if err != nil {
		return nil
	}
	defer rows.Close()
	cts, err := rows.ColumnTypes()
	if err != nil || len(cts) != 1 {
		return nil
	}
	return typeOf(cts[0])
}

// insertTargets reads the table and columns an INSERT ... VALUES writes,
// as the note prints them: the columns listed, or every column of the
// table in order when none are.
func (s *statement) insertTargets(t text) (*relation, []string) {
	for i, tok := range t.toks {
		if tok.kind != tkWord || tok.text != "into" || t.at(i+1).kind != tkIdent {
			continue
		}
		parts := t.at(i + 1).parts
		schema := ""
		if len(parts) >= 2 {
			schema = parts[len(parts)-2]
		}
		rel, ok := s.a.table(schema, parts[len(parts)-1])
		if !ok {
			return nil, nil
		}
		var targets []string
		if t.at(i+2).kind == tkLParen {
			for _, m := range t.list(i + 2) {
				if m[1]-m[0] == 1 && t.at(m[0]).kind == tkIdent {
					parts := t.at(m[0]).parts
					targets = append(targets, parts[len(parts)-1])
				}
			}
		} else {
			for _, col := range s.a.catalog[rel] {
				targets = append(targets, col.name)
			}
		}
		return &rel, targets
	}
	return nil, nil
}
