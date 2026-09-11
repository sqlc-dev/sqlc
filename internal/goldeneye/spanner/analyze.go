package spanner

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"cloud.google.com/go/spanner/apiv1/spannerpb"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// The analyze cases are checked against a live server, which compiles
// each query in PLAN mode: nothing runs, and the server reports the name
// and type of each result column, the type of each parameter the query
// leaves undeclared, and the query plan. The metadata types a value the
// way the wire does — STRING, ARRAY<INT64> — without the length a
// declaration gives a STRING(10) or whether the column can be NULL, so a
// result column the plan reads from a table column is described from the
// case's INFORMATION_SCHEMA instead, and a parameter the plan compares
// with or assigns to a table column is described as that column.

// placeholder is one parameter of a query as sqlc numbers them: each
// @name or sqlc.arg name once, at its first appearance, which is how
// Spanner numbers its named parameters too.
type placeholder struct {
	Number int
	Name   string
}

var identRe = regexp.MustCompile(`^@([A-Za-z_][A-Za-z0-9_]*)`)

// bind rewrites the query so that every parameter is a Spanner @name, and
// lists the parameters in sqlc's order. The query's own @name references
// are kept, since that is how sqlc's GoogleSQL queries name their
// parameters; a ? or a sqlc.arg becomes one.
func bind(query string) (string, []placeholder) {
	sql := endtoend.Rewrite(query, func(name, _ string) string {
		if name == "" {
			name = "p"
		}
		return "@" + name
	})
	var phs []placeholder
	numbers := map[string]int{}
	i := 0
	for i < len(sql) {
		c := sql[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			i = skipQuoted(sql, i)
		case strings.HasPrefix(sql[i:], "--") || strings.HasPrefix(sql[i:], "#"):
			end := strings.IndexByte(sql[i:], '\n')
			if end < 0 {
				i = len(sql)
			} else {
				i += end
			}
		case strings.HasPrefix(sql[i:], "/*"):
			end := strings.Index(sql[i:], "*/")
			if end < 0 {
				i = len(sql)
			} else {
				i += end + 2
			}
		case c == '@' && identRe.MatchString(sql[i:]):
			m := identRe.FindStringSubmatch(sql[i:])
			if _, ok := numbers[m[1]]; !ok {
				numbers[m[1]] = len(phs) + 1
				phs = append(phs, placeholder{Number: len(phs) + 1, Name: m[1]})
			}
			i += len(m[0])
		default:
			i++
		}
	}
	return sql, phs
}

// skipQuoted returns the index just past the quoted token starting at i,
// honouring backslash escapes.
func skipQuoted(s string, i int) int {
	q := s[i]
	j := i + 1
	for j < len(s) {
		switch {
		case s[j] == '\\' && j+1 < len(s):
			j += 2
		case s[j] == q:
			return j + 1
		default:
			j++
		}
	}
	return len(s)
}

// splitStatements splits a script on the semicolons outside strings and
// comments, which is how a schema's DDL is handed to CREATE DATABASE and a
// fixture's DML to a transaction.
func splitStatements(src string) []string {
	var stmts []string
	flush := func(s string) {
		if s = strings.TrimSpace(s); s != "" {
			stmts = append(stmts, s)
		}
	}
	start, i := 0, 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			i = skipQuoted(src, i)
		case strings.HasPrefix(src[i:], "--") || strings.HasPrefix(src[i:], "#"):
			end := strings.IndexByte(src[i:], '\n')
			if end < 0 {
				i = len(src)
			} else {
				i += end
			}
		case strings.HasPrefix(src[i:], "/*"):
			end := strings.Index(src[i:], "*/")
			if end < 0 {
				i = len(src)
			} else {
				i += end + 2
			}
		case c == ';':
			flush(src[start:i])
			start = i + 1
			i++
		default:
			i++
		}
	}
	flush(src[start:])
	return stmts
}

// column is what INFORMATION_SCHEMA says about a column of the case's
// database: its type as the schema declared it, and whether it can be
// NULL.
type column struct {
	name     string
	typ      *analysis.TypeExpr
	nullable bool
}

// analyzer holds the session a case is compiled in and the catalog of the
// case's database, keyed by table name in the database's own schema.
type analyzer struct {
	s       *server
	session string
	catalog map[string][]column
	keys    map[string][]string // the primary key columns of each table
}

var dbNameRe = regexp.MustCompile(`[^a-z0-9_]+`)

// Analyze creates a database of the case's own from its schema, writes its
// fixture there, compiles its queries and returns what Spanner reports in
// the JSON shape sqlc analyze prints.
func Analyze(ctx context.Context, endpoint string, c endtoend.Case) ([]byte, error) {
	s, err := open(ctx, endpoint)
	if err != nil {
		return nil, err
	}
	defer s.Close()

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

	// A database id is lower-case letters, digits and underscores, at
	// most 30 characters.
	name := "goldeneye_" + dbNameRe.ReplaceAllString(strings.ToLower(c.Name), "_")
	if len(name) > 30 {
		name = name[:30]
	}
	name = strings.TrimRight(name, "_")
	db := Instance + "/databases/" + name
	if err := s.dropDatabase(ctx, db); err != nil {
		return nil, err
	}
	if _, err := s.createDatabase(ctx, name, splitStatements(string(schema))); err != nil {
		return nil, fmt.Errorf("loading %s: %w", c.Schema, err)
	}
	defer s.dropDatabase(context.WithoutCancel(ctx), db)
	session, err := s.session(ctx, db)
	if err != nil {
		return nil, err
	}
	if stmts := splitStatements(string(fixture)); len(stmts) > 0 {
		if err := s.write(ctx, session, stmts); err != nil {
			return nil, fmt.Errorf("loading %s: %w", c.Fixture, err)
		}
	}

	a := &analyzer{s: s, session: session}
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

// Check compares what Spanner reports for a case with the output the
// case committed, returning a diff when they differ.
func Check(ctx context.Context, endpoint string, c endtoend.Case) (string, error) {
	got, err := Analyze(ctx, endpoint, c)
	if err != nil {
		return "", err
	}
	return c.Compare(got)
}

// write runs DML statements in a read-write transaction and commits it.
func (s *server) write(ctx context.Context, session string, stmts []string) error {
	tx, err := s.data.BeginTransaction(ctx, &spannerpb.BeginTransactionRequest{
		Session: session,
		Options: &spannerpb.TransactionOptions{Mode: &spannerpb.TransactionOptions_ReadWrite_{ReadWrite: &spannerpb.TransactionOptions_ReadWrite{}}},
	})
	if err != nil {
		return err
	}
	req := &spannerpb.ExecuteBatchDmlRequest{
		Session:     session,
		Transaction: &spannerpb.TransactionSelector{Selector: &spannerpb.TransactionSelector_Id{Id: tx.Id}},
		Seqno:       1,
	}
	for _, stmt := range stmts {
		req.Statements = append(req.Statements, &spannerpb.ExecuteBatchDmlRequest_Statement{Sql: stmt})
	}
	resp, err := s.data.ExecuteBatchDml(ctx, req)
	if err != nil {
		return err
	}
	if resp.Status != nil && resp.Status.Code != 0 {
		return fmt.Errorf("%s", resp.Status.Message)
	}
	_, err = s.data.Commit(ctx, &spannerpb.CommitRequest{
		Session:     session,
		Transaction: &spannerpb.CommitRequest_TransactionId{TransactionId: tx.Id},
	})
	return err
}

const catalogQuery = `
SELECT TABLE_NAME, COLUMN_NAME, SPANNER_TYPE, IS_NULLABLE
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_SCHEMA = ''
ORDER BY TABLE_NAME, ORDINAL_POSITION`

const keysQuery = `
SELECT TABLE_NAME, COLUMN_NAME
FROM INFORMATION_SCHEMA.INDEX_COLUMNS
WHERE TABLE_SCHEMA = '' AND INDEX_TYPE = 'PRIMARY_KEY'
ORDER BY TABLE_NAME, ORDINAL_POSITION`

// readCatalog reads the tables the schema created, and their keys.
func (a *analyzer) readCatalog(ctx context.Context) error {
	rows, err := a.s.query(ctx, a.session, catalogQuery)
	if err != nil {
		return err
	}
	a.catalog = map[string][]column{}
	for _, row := range rows {
		table := row[0].GetStringValue()
		a.catalog[table] = append(a.catalog[table], column{
			name:     row[1].GetStringValue(),
			typ:      parseType(row[2].GetStringValue()),
			nullable: row[3].GetStringValue() == "YES",
		})
	}
	rows, err = a.s.query(ctx, a.session, keysQuery)
	if err != nil {
		return err
	}
	a.keys = map[string][]string{}
	for _, row := range rows {
		table := row[0].GetStringValue()
		a.keys[table] = append(a.keys[table], row[1].GetStringValue())
	}
	return nil
}

// tableName returns the catalog's spelling of a table's name, which
// Spanner matches in any case.
func (a *analyzer) tableName(table string) string {
	for t := range a.catalog {
		if strings.EqualFold(t, table) {
			return t
		}
	}
	return table
}

// lookup finds a column of a table of the case's database. Spanner
// matches table and column names in any case.
func (a *analyzer) lookup(table, name string) (column, bool) {
	for t, cols := range a.catalog {
		if !strings.EqualFold(t, table) {
			continue
		}
		for _, col := range cols {
			if strings.EqualFold(col.name, name) {
				return col, true
			}
		}
	}
	return column{}, false
}

// parseType reads a type spelled the way SPANNER_TYPE and a declaration
// spell it — STRING(10), ARRAY<STRING(MAX)>, STRUCT<a INT64, b STRING> —
// into an expression, in lower case, with MAX an identifier.
func parseType(s string) *analysis.TypeExpr {
	s = strings.TrimSpace(s)
	lower := strings.ToLower(s)
	if element, ok := strings.CutPrefix(lower, "array<"); ok && strings.HasSuffix(element, ">") {
		return &analysis.TypeExpr{Name: "array", Args: []analysis.TypeArg{{Type: parseType(s[6 : len(s)-1])}}}
	}
	if fields, ok := strings.CutPrefix(lower, "struct<"); ok && strings.HasSuffix(fields, ">") {
		t := &analysis.TypeExpr{Name: "struct"}
		for _, f := range splitTop(s[7:len(s)-1], ',') {
			f = strings.TrimSpace(f)
			label, typ := "", f
			if i := strings.IndexAny(f, " \t"); i > 0 && !strings.ContainsAny(f[:i], "<(") {
				label, typ = f[:i], strings.TrimSpace(f[i+1:])
			}
			t.Args = append(t.Args, analysis.TypeArg{Label: label, Type: parseType(typ)})
		}
		return t
	}
	name, args := lower, ""
	if open := strings.IndexByte(lower, '('); open >= 0 && strings.HasSuffix(lower, ")") {
		name, args = strings.TrimSpace(lower[:open]), lower[open+1:len(lower)-1]
	}
	t := &analysis.TypeExpr{Name: name}
	if args == "" {
		return t
	}
	for _, arg := range strings.Split(args, ",") {
		arg = strings.TrimSpace(arg)
		if n, err := strconv.ParseInt(arg, 10, 64); err == nil {
			t.Args = append(t.Args, analysis.TypeArg{Int: &n})
		} else {
			ident := arg
			t.Args = append(t.Args, analysis.TypeArg{Ident: &ident})
		}
	}
	return t
}

// splitTop splits on a separator outside angle brackets and parentheses.
func splitTop(s string, sep byte) []string {
	var out []string
	depth, start := 0, 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '<', '(':
			depth++
		case '>', ')':
			depth--
		case sep:
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	return append(out, s[start:])
}

// typeOf reads a type the way the wire spells it, which names the family
// and, for an array or a struct, what it holds, and nothing of a length.
func typeOf(t *spannerpb.Type) *analysis.TypeExpr {
	if t == nil {
		return nil
	}
	switch t.Code {
	case spannerpb.TypeCode_ARRAY:
		return &analysis.TypeExpr{Name: "array", Args: []analysis.TypeArg{{Type: typeOf(t.ArrayElementType)}}}
	case spannerpb.TypeCode_STRUCT:
		out := &analysis.TypeExpr{Name: "struct"}
		if t.StructType != nil {
			for _, f := range t.StructType.Fields {
				out.Args = append(out.Args, analysis.TypeArg{Label: f.Name, Type: typeOf(f.Type)})
			}
		}
		return out
	case spannerpb.TypeCode_PROTO, spannerpb.TypeCode_ENUM:
		return &analysis.TypeExpr{Name: strings.ToLower(t.ProtoTypeFqn)}
	}
	return &analysis.TypeExpr{Name: strings.ToLower(t.Code.String())}
}

// withNullable copies a type with its nullability set.
func withNullable(t *analysis.TypeExpr, nullable bool) *analysis.TypeExpr {
	if t == nil {
		return nil
	}
	out := *t
	out.Nullable = nullable
	return &out
}

// isDML reports whether a statement writes, and so has to be compiled in a
// read-write transaction, which is begun for it and never committed.
func isDML(sql string) bool {
	head := strings.ToLower(strings.TrimSpace(sql))
	for _, kw := range []string{"insert", "update", "delete"} {
		if strings.HasPrefix(head, kw) {
			return true
		}
	}
	return false
}

// compile compiles a statement in PLAN mode and returns what the server
// reports about it.
func (a *analyzer) compile(ctx context.Context, sql string) (*spannerpb.ResultSet, error) {
	req := &spannerpb.ExecuteSqlRequest{
		Session:   a.session,
		Sql:       sql,
		QueryMode: spannerpb.ExecuteSqlRequest_PLAN,
	}
	if isDML(sql) {
		req.Transaction = &spannerpb.TransactionSelector{Selector: &spannerpb.TransactionSelector_Begin{
			Begin: &spannerpb.TransactionOptions{Mode: &spannerpb.TransactionOptions_ReadWrite_{ReadWrite: &spannerpb.TransactionOptions_ReadWrite{}}},
		}}
		req.Seqno = 1
	}
	rs, err := a.s.data.ExecuteSql(ctx, req)
	if err != nil {
		return nil, err
	}
	if rs.Metadata != nil && rs.Metadata.Transaction != nil && len(rs.Metadata.Transaction.Id) > 0 {
		a.s.data.Rollback(context.WithoutCancel(ctx), &spannerpb.RollbackRequest{Session: a.session, TransactionId: rs.Metadata.Transaction.Id})
	}
	return rs, nil
}

// analyzeQuery compiles one query.
func (a *analyzer) analyzeQuery(ctx context.Context, q endtoend.Query) (analysis.Query, error) {
	sql, phs := bind(q.SQL)
	aq := analysis.Query{
		Name:    q.Name,
		Cmd:     q.Cmd,
		Columns: []analysis.Column{},
		Params:  []analysis.Param{},
	}
	rs, err := a.compile(ctx, sql)
	if err != nil {
		return analysis.Query{}, err
	}
	p := newPlan(rs.Stats.GetQueryPlan())
	partners := p.partners()
	describe := func(o origin) (analysis.Column, bool) {
		if o.param != "" {
			o = partners[o.param]
		}
		if o.table == "" {
			return analysis.Column{}, false
		}
		col, ok := a.lookup(o.table, o.column)
		if !ok {
			return analysis.Column{}, false
		}
		return analysis.Column{Name: col.name, Type: withNullable(col.typ, col.nullable), Table: a.tableName(o.table)}, true
	}

	var fields []*spannerpb.StructType_Field
	if rs.Metadata != nil && rs.Metadata.RowType != nil {
		fields = rs.Metadata.RowType.Fields
	}
	outputs := p.outputs()
	table, operation := p.mutation()
	if table != "" {
		// The values a DML plan writes come before the columns it
		// returns: each is a parameter standing for the column it is
		// written to.
		written := a.written(sql, table, operation)
		for i, w := range written {
			if i >= len(outputs)-len(fields) {
				break
			}
			if o := p.resolve(outputs[i], map[int32]bool{}); o.param != "" {
				if _, ok := partners[o.param]; !ok {
					partners[o.param] = origin{table: table, column: w}
				}
			}
		}
		outputs = outputs[max(0, len(outputs)-len(fields)):]
	}
	for i, f := range fields {
		ac := analysis.Column{Name: f.Name, Type: typeOf(f.Type)}
		switch {
		case table != "":
			// A THEN RETURN column is the table's column of that name.
			if col, ok := a.lookup(table, f.Name); ok {
				ac.Type = withNullable(col.typ, col.nullable)
				ac.Table = a.tableName(table)
			}
		case i < len(outputs):
			if col, ok := describe(p.resolve(outputs[i], map[int32]bool{})); ok {
				ac.Type, ac.Table = col.Type, col.Table
			}
		}
		aq.Columns = append(aq.Columns, ac)
	}

	declared := map[string]*spannerpb.Type{}
	if rs.Metadata != nil && rs.Metadata.UndeclaredParameters != nil {
		for _, f := range rs.Metadata.UndeclaredParameters.Fields {
			declared[f.Name] = f.Type
		}
	}
	for _, ph := range phs {
		ac, ok := describe(origin{param: ph.Name})
		if !ok {
			ac = analysis.Column{Type: typeOf(declared[ph.Name])}
		}
		aq.Params = append(aq.Params, analysis.Param{Number: ph.Number, Column: ac})
	}
	return aq, nil
}

var (
	insertRe = regexp.MustCompile("(?is)^insert\\s+(?:or\\s+\\w+\\s+)?into\\s+[\\w.`]+\\s*(?:\\(([^)]*)\\))?")
	updateRe = regexp.MustCompile("(?is)^update\\s+[\\w.`]+(?:\\s+(?:as\\s+)?\\w+)?\\s+set\\s+(.*?)\\s+where\\b")
)

// written lists the columns a DML statement writes, in the order the plan
// lists their values: the table's key columns, then for an UPDATE the
// columns it sets and for an INSERT the columns it inserts, which are the
// statement's column list or every column of the table.
func (a *analyzer) written(sql, table, operation string) []string {
	switch operation {
	case "INSERT":
		m := insertRe.FindStringSubmatch(sql)
		if m == nil {
			return nil
		}
		if strings.TrimSpace(m[1]) == "" {
			var cols []string
			for _, col := range a.catalog[a.tableName(table)] {
				cols = append(cols, col.name)
			}
			return cols
		}
		var cols []string
		for _, c := range strings.Split(m[1], ",") {
			cols = append(cols, strings.Trim(strings.TrimSpace(c), "`"))
		}
		return cols
	case "UPDATE":
		cols := append([]string(nil), a.keys[a.tableName(table)]...)
		if m := updateRe.FindStringSubmatch(sql); m != nil {
			for _, assignment := range splitTop(m[1], ',') {
				target, _, _ := strings.Cut(assignment, "=")
				target = strings.TrimSpace(target)
				if i := strings.LastIndexByte(target, '.'); i >= 0 {
					target = target[i+1:]
				}
				cols = append(cols, strings.Trim(target, "`"))
			}
		}
		return cols
	case "DELETE":
		return a.keys[a.tableName(table)]
	}
	return nil
}
