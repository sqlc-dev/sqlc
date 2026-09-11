package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sqlc-dev/sqlc/internal/goldeneye/analysis"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/dialect"
	"github.com/sqlc-dev/sqlc/internal/goldeneye/endtoend"
)

// The analyze cases are checked against a live server, which is asked
// three things about each query, none of which runs it. What a driver
// would see: sys.dm_exec_describe_first_result_set describes each result
// column — its name, its type spelled the way a declaration spells it,
// whether it can be NULL, and which table column it is read from. What
// each parameter would be: sp_describe_undeclared_parameters says what
// type the server would give each parameter the query leaves undeclared.
// And what each parameter stands in for: the estimated showplan, compiled
// with the parameters declared as those types, says which column each is
// compared with or assigned to, and a parameter with such a partner is
// described as that column, from the catalog of the case's database.
//
// Two things the describing function keeps to itself: a json or vector
// column is described by the nvarchar(max) it is sent to a driver as, so
// a column read from a table is typed from the catalog instead, and a
// type SQL Server calls timestamp is written as the rowversion the
// dialect names it.

// placeholder is one parameter of a query as sqlc numbers them: each ?
// in turn, and each @name or sqlc.arg name once, at its first appearance.
// Each appearance becomes a variable of its own, since the server
// describes an undeclared parameter only when it is used once.
type placeholder struct {
	Number int
	Name   string
	Vars   []string // the variables its appearances became, without the @
}

var identRe = regexp.MustCompile(`^@([A-Za-z_][A-Za-z0-9_]*)`)

// bind rewrites the query so that every parameter is a variable used once,
// and lists the parameters in sqlc's order. The query's own @name
// references are kept, since that is how sqlc's SQL Server queries name
// their parameters; a ? or a sqlc.arg becomes one.
func bind(query string) (string, []placeholder) {
	sql := endtoend.Rewrite(query, func(name, _ string) string {
		if name == "" {
			name = "p"
		}
		return "@" + name
	})
	var (
		phs     []placeholder
		numbers = map[string]int{}
		out     strings.Builder
	)
	i := 0
	for i < len(sql) {
		c := sql[i]
		switch {
		case c == '\'' || c == '"' || c == '[':
			end := skipQuoted(sql, i)
			out.WriteString(sql[i:end])
			i = end
		case strings.HasPrefix(sql[i:], "--"):
			end := strings.IndexByte(sql[i:], '\n')
			if end < 0 {
				end = len(sql)
			} else {
				end += i
			}
			out.WriteString(sql[i:end])
			i = end
		case strings.HasPrefix(sql[i:], "/*"):
			end := strings.Index(sql[i:], "*/")
			if end < 0 {
				end = len(sql)
			} else {
				end += i + 2
			}
			out.WriteString(sql[i:end])
			i = end
		case c == '@' && (i == 0 || !isWordByte(sql[i-1])) && identRe.MatchString(sql[i:]):
			m := identRe.FindStringSubmatch(sql[i:])
			name := m[1]
			n, ok := numbers[name]
			if !ok {
				n = len(phs) + 1
				numbers[name] = n
				phs = append(phs, placeholder{Number: n, Name: name})
			}
			ph := &phs[n-1]
			v := name
			if k := len(ph.Vars); k > 0 {
				v = fmt.Sprintf("%s__%d", name, k+1)
			}
			ph.Vars = append(ph.Vars, v)
			out.WriteString("@" + v)
			i += len(m[0])
		default:
			out.WriteByte(c)
			i++
		}
	}
	return out.String(), phs
}

func isWordByte(c byte) bool {
	return c == '_' || c == '@' || c >= '0' && c <= '9' || c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

// skipQuoted returns the index just past the quoted token starting at i: a
// string, a double-quoted identifier or a bracketed one, whose closing
// delimiter is escaped by doubling it.
func skipQuoted(s string, i int) int {
	open := s[i]
	close := open
	if open == '[' {
		close = ']'
	}
	j := i + 1
	for j < len(s) {
		switch {
		case s[j] == close && j+1 < len(s) && s[j+1] == close:
			j += 2
		case s[j] == close:
			return j + 1
		default:
			j++
		}
	}
	return len(s)
}

// splitStatements splits a script into the statements it is made of, on
// the semicolons outside strings, brackets and comments and on the GO
// lines a T-SQL script separates batches with. A CREATE TYPE has to be
// its own batch before a table can use the type, so a schema is loaded
// one statement at a time.
func splitStatements(src string) []string {
	var stmts []string
	flush := func(s string) {
		if s = strings.TrimSpace(s); s != "" && !strings.EqualFold(s, "go") {
			stmts = append(stmts, s)
		}
	}
	start := 0
	i := 0
	for i < len(src) {
		c := src[i]
		switch {
		case c == '\'' || c == '"' || c == '[':
			i = skipQuoted(src, i)
		case strings.HasPrefix(src[i:], "--"):
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
		case c == '\n' && isGoLine(src, i+1):
			flush(src[start:i])
			i++
			for i < len(src) && src[i] != '\n' {
				i++
			}
			start = i
		default:
			i++
		}
	}
	flush(src[start:])
	return stmts
}

// isGoLine reports whether the line starting at i is a GO batch separator.
func isGoLine(src string, i int) bool {
	end := strings.IndexByte(src[i:], '\n')
	if end < 0 {
		end = len(src)
	} else {
		end += i
	}
	return strings.EqualFold(strings.TrimSpace(src[i:end]), "go")
}

// column is what the catalog says about a column of the case's database.
type column struct {
	name     string
	typ      *analysis.TypeExpr
	nullable bool
}

// relation is a table or view of the case's database, by schema and name,
// in lower case.
type relation struct {
	schema, name string
}

// analyzer holds the session a case is described in and the catalog of
// the case's database.
type analyzer struct {
	conn    *sql.Conn
	catalog map[relation][]column
	// userTypes are the alias types the schema created, by name in lower
	// case, which the server reports beside the system type they stand on
	// and the dialect reports by their own name.
	userTypes map[string]bool
	// canonical names the type the dialect reports a spelling as, for each
	// alias types.jsonl lists: SQL Server spells numeric and timestamp in
	// its catalog, which the dialect names decimal and rowversion.
	canonical map[string]string
}

var dbNameRe = regexp.MustCompile(`[^A-Za-z0-9_]+`)

// Analyze loads a case's schema and fixture into a database of their own
// on the server, describes its queries there and returns what SQL Server
// reports in the JSON shape sqlc analyze prints.
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

	db := quote("goldeneye_" + dbNameRe.ReplaceAllString(c.Name, "_"))
	for _, stmt := range []string{
		"USE master",
		"DROP DATABASE IF EXISTS " + db,
		"CREATE DATABASE " + db,
		"USE " + db,
	} {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return nil, err
		}
	}
	defer conn.ExecContext(context.WithoutCancel(ctx), "USE master; DROP DATABASE IF EXISTS "+db)
	for _, stmt := range splitStatements(string(schema)) {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return nil, fmt.Errorf("loading %s: %w", c.Schema, err)
		}
	}
	for _, stmt := range splitStatements(string(fixture)) {
		if _, err := conn.ExecContext(ctx, stmt); err != nil {
			return nil, fmt.Errorf("loading %s: %w", c.Fixture, err)
		}
	}

	a := &analyzer{conn: conn}
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

// readAliases reads the aliases the hand-written types.jsonl gives each
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

// Check compares what SQL Server reports for a case with the output the
// case committed, returning a diff when they differ.
func Check(ctx context.Context, dsn string, c endtoend.Case) (string, error) {
	got, err := Analyze(ctx, dsn, c)
	if err != nil {
		return "", err
	}
	return c.Compare(got)
}

// tablesQuery lists the tables and views of the current database, and the
// columns of each: the type as the schema declared it, whether that is a
// type the schema created, and the dimensions of a vector, which is all
// the describing function does not say.
const tablesQuery = `
SELECT s.name, o.name, c.name, t.name, t.is_user_defined, c.vector_dimensions
FROM sys.columns c
JOIN sys.objects o ON o.object_id = c.object_id
JOIN sys.schemas s ON s.schema_id = o.schema_id
JOIN sys.types t ON t.user_type_id = c.user_type_id
WHERE o.type IN ('U', 'V')
ORDER BY s.name, o.name, c.column_id`

// readCatalog reads the tables the schema created. The server spells each
// column's type by describing a SELECT * from the table, and the catalog
// says which columns that spelling is not the declared type of.
func (a *analyzer) readCatalog(ctx context.Context) error {
	a.catalog = map[relation][]column{}
	a.userTypes = map[string]bool{}
	rows, err := a.conn.QueryContext(ctx, tablesQuery)
	if err != nil {
		return err
	}
	type declared struct {
		typeName    string
		userDefined bool
		dimensions  sql.NullInt64
	}
	decls := map[relation][]declared{}
	var order []relation
	for rows.Next() {
		var (
			schema, table, col string
			d                  declared
		)
		if err := rows.Scan(&schema, &table, &col, &d.typeName, &d.userDefined, &d.dimensions); err != nil {
			rows.Close()
			return err
		}
		rel := relation{strings.ToLower(schema), strings.ToLower(table)}
		if _, ok := decls[rel]; !ok {
			order = append(order, rel)
		}
		decls[rel] = append(decls[rel], d)
		if d.userDefined {
			a.userTypes[strings.ToLower(d.typeName)] = true
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, rel := range order {
		described, err := a.describe(ctx, "SELECT * FROM "+quote(rel.schema)+"."+quote(rel.name))
		if err != nil {
			return fmt.Errorf("%s.%s: %w", rel.schema, rel.name, err)
		}
		if len(described) != len(decls[rel]) {
			return fmt.Errorf("%s.%s: the catalog lists %d columns and the server describes %d", rel.schema, rel.name, len(decls[rel]), len(described))
		}
		var cols []column
		for i, dc := range described {
			d := decls[rel][i]
			t := dc.typ
			switch {
			case d.userDefined:
				t = &analysis.TypeExpr{Name: strings.ToLower(d.typeName)}
			case strings.EqualFold(d.typeName, "json"):
				t = &analysis.TypeExpr{Name: "json"}
			case strings.EqualFold(d.typeName, "vector") && d.dimensions.Valid:
				n := d.dimensions.Int64
				t = &analysis.TypeExpr{Name: "vector", Args: []analysis.TypeArg{{Int: &n}}}
			}
			cols = append(cols, column{name: strings.ToLower(dc.name), typ: t, nullable: dc.nullable})
		}
		a.catalog[rel] = cols
	}
	return nil
}

// described is one result column as sys.dm_exec_describe_first_result_set
// reports it.
type described struct {
	name     string
	typ      *analysis.TypeExpr
	nullable bool
	source   *columnRef // the table column it is read from, or nil
}

const describeColumnsQuery = `
SELECT column_ordinal, name, system_type_name, user_type_name, is_nullable,
  source_schema, source_table, source_column, is_computed_column, error_message
FROM sys.dm_exec_describe_first_result_set(@p1, NULL, 1)
WHERE is_hidden = 0
ORDER BY column_ordinal`

// describe asks the server what the first result set of a statement
// looks like. A statement that returns no rows describes as no columns; one
// the server cannot compile describes as its error.
func (a *analyzer) describe(ctx context.Context, query string) ([]described, error) {
	rows, err := a.conn.QueryContext(ctx, describeColumnsQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []described
	for rows.Next() {
		var (
			ordinal                        int
			name, sysType, userType        sql.NullString
			nullable                       sql.NullBool
			srcSchema, srcTable, srcColumn sql.NullString
			computed                       sql.NullBool
			errMsg                         sql.NullString
		)
		if err := rows.Scan(&ordinal, &name, &sysType, &userType, &nullable, &srcSchema, &srcTable, &srcColumn, &computed, &errMsg); err != nil {
			return nil, err
		}
		if errMsg.Valid {
			return nil, fmt.Errorf("cannot be described: %s", errMsg.String)
		}
		d := described{
			name:     name.String,
			typ:      a.typeOf(sysType.String, userType),
			nullable: nullable.Bool,
		}
		// A computed column names the column it is computed from, as the
		// one an update through the result set would write; a column is
		// read from a table only when it is that column.
		if srcTable.Valid && !computed.Bool {
			d.source = &columnRef{schema: srcSchema.String, table: srcTable.String, column: srcColumn.String}
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// typeOf reads a type the server spelled, which is the system type unless
// the server also names a type the schema created, which the dialect
// reports by that name and no argument.
func (a *analyzer) typeOf(systemType string, userType sql.NullString) *analysis.TypeExpr {
	if userType.Valid && a.userTypes[strings.ToLower(userType.String)] {
		return &analysis.TypeExpr{Name: strings.ToLower(userType.String)}
	}
	return parseType(systemType, a.canonical)
}

// parseType reads a type spelled the way a declaration spells it —
// nvarchar(max), decimal(10,2), datetime2(7) — into an expression, named
// the way the dialect names it when the spelling is one of the aliases
// the dialect lists: numeric is decimal, timestamp is rowversion.
func parseType(s string, canonical map[string]string) *analysis.TypeExpr {
	s = strings.TrimSpace(strings.ToLower(s))
	name, args := s, ""
	if open := strings.IndexByte(s, '('); open >= 0 && strings.HasSuffix(s, ")") {
		name, args = strings.TrimSpace(s[:open]), s[open+1:len(s)-1]
	}
	if c, ok := canonical[name]; ok {
		name = c
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

// withNullable copies a type with its nullability set.
func withNullable(t *analysis.TypeExpr, nullable bool) *analysis.TypeExpr {
	if t == nil {
		return nil
	}
	out := *t
	out.Nullable = nullable
	return &out
}

// lookup finds a column of the case's database in the catalog.
func (a *analyzer) lookup(ref columnRef) (column, bool) {
	cols, ok := a.catalog[relation{strings.ToLower(ref.schema), strings.ToLower(ref.table)}]
	if !ok {
		return column{}, false
	}
	for _, col := range cols {
		if strings.EqualFold(col.name, ref.column) {
			return col, true
		}
	}
	return column{}, false
}

// parameter is what sp_describe_undeclared_parameters says about one
// variable the statement leaves undeclared.
type parameter struct {
	systemType string
	userType   sql.NullString
}

// undeclaredParameters asks the server what type it would give each
// parameter, keyed by the variable's name without the @.
func (a *analyzer) undeclaredParameters(ctx context.Context, query string) (map[string]parameter, error) {
	rows, err := a.conn.QueryContext(ctx, "EXEC sp_describe_undeclared_parameters @p1", query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	index := map[string]int{}
	for i, c := range cols {
		index[c] = i
	}
	params := map[string]parameter{}
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		str := func(col string) sql.NullString {
			i, ok := index[col]
			if !ok || vals[i] == nil {
				return sql.NullString{}
			}
			return sql.NullString{String: fmt.Sprint(vals[i]), Valid: true}
		}
		name := strings.TrimPrefix(str("name").String, "@")
		params[name] = parameter{
			systemType: str("suggested_system_type_name").String,
			userType:   str("suggested_user_type_name"),
		}
	}
	return params, rows.Err()
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

	described, err := a.describe(ctx, sql)
	if err != nil {
		return analysis.Query{}, err
	}
	for _, d := range described {
		ac := analysis.Column{Name: strings.ToLower(d.name), Type: withNullable(d.typ, d.nullable)}
		if d.source != nil {
			if col, ok := a.lookup(*d.source); ok {
				ac.Type = withNullable(col.typ, d.nullable)
			}
			ac.Table = strings.ToLower(d.source.table)
		}
		aq.Columns = append(aq.Columns, ac)
	}

	if len(phs) == 0 {
		return aq, nil
	}
	params, err := a.undeclaredParameters(ctx, sql)
	if err != nil {
		return analysis.Query{}, err
	}
	var decl []string
	for _, ph := range phs {
		for _, v := range ph.Vars {
			p, ok := params[v]
			if !ok {
				return analysis.Query{}, fmt.Errorf("the server did not describe parameter @%s", v)
			}
			decl = append(decl, "@"+v+" "+p.systemType)
		}
	}
	partners, err := a.partners(ctx, "DECLARE "+strings.Join(decl, ", ")+";\n"+sql)
	if err != nil {
		return analysis.Query{}, err
	}
	for _, ph := range phs {
		ac := analysis.Column{}
		found := false
		for _, v := range ph.Vars {
			ref, ok := partners[v]
			if !ok {
				continue
			}
			col, ok := a.lookup(ref)
			if !ok {
				continue
			}
			ac = analysis.Column{Name: col.name, Type: withNullable(col.typ, col.nullable), Table: strings.ToLower(ref.table)}
			found = true
			break
		}
		if !found {
			p := params[ph.Vars[0]]
			ac.Type = a.typeOf(p.systemType, p.userType)
		}
		aq.Params = append(aq.Params, analysis.Param{Number: ph.Number, Column: ac})
	}
	return aq, nil
}
