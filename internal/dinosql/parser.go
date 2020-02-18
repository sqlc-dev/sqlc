package dinosql

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/kyleconroy/sqlc/internal/catalog"
	core "github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/postgres"
	"github.com/kyleconroy/sqlc/internal/postgresql/ast"

	"github.com/davecgh/go-spew/spew"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func keepSpew() {
	spew.Dump("hello world")
}

type FileErr struct {
	Filename string
	Line     int
	Column   int
	Err      error
}

type ParserErr struct {
	Errs []FileErr
}

func (e *ParserErr) Add(filename, source string, loc int, err error) {
	line := 1
	column := 1
	if lerr, ok := err.(core.Error); ok {
		if lerr.Location != 0 {
			loc = lerr.Location
		}
	}
	if source != "" && loc != 0 {
		line, column = lineno(source, loc)
	}
	e.Errs = append(e.Errs, FileErr{filename, line, column, err})
}

func NewParserErr() *ParserErr {
	return &ParserErr{}
}

func (e *ParserErr) Error() string {
	return fmt.Sprintf("multiple errors: %d errors", len(e.Errs))
}

func ReadSQLFiles(path string) ([]string, error) {
	f, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("path %s does not exist", path)
	}

	var files []string
	if f.IsDir() {
		listing, err := ioutil.ReadDir(path)
		if err != nil {
			return nil, err
		}
		for _, f := range listing {
			files = append(files, filepath.Join(path, f.Name()))
		}
	} else {
		files = append(files, path)
	}

	var sql []string
	for _, filename := range files {
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}
		if strings.HasPrefix(filepath.Base(filename), ".") {
			continue
		}
		// Remove golang-migrate rollback files.
		if strings.HasSuffix(filename, ".down.sql") {
			continue
		}
		sql = append(sql, filename)
	}
	return sql, nil
}

func ParseCatalog(schema string) (core.Catalog, error) {
	files, err := ReadSQLFiles(schema)
	if err != nil {
		return core.Catalog{}, err
	}

	merr := NewParserErr()
	c := core.NewCatalog()
	for _, filename := range files {
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		contents := RemoveRollbackStatements(string(blob))
		tree, err := pg.Parse(contents)
		if err != nil {
			merr.Add(filename, contents, 0, err)
			continue
		}
		for _, stmt := range tree.Statements {
			if err := validateFuncCall(&c, stmt); err != nil {
				merr.Add(filename, contents, location(stmt), err)
				continue
			}
			if err := catalog.Update(&c, stmt); err != nil {
				merr.Add(filename, contents, location(stmt), err)
				continue
			}
		}
	}

	// The pg_temp schema is scoped to the current session. Remove it from the
	// catalog so that other queries can not read from it.
	delete(c.Schemas, "pg_temp")

	if len(merr.Errs) > 0 {
		return c, merr
	}
	return c, nil
}

func updateCatalog(c *core.Catalog, tree pg.ParsetreeList) error {
	for _, stmt := range tree.Statements {
		if err := validateFuncCall(c, stmt); err != nil {
			return err
		}
		if err := catalog.Update(c, stmt); err != nil {
			return err
		}
	}
	return nil
}

func join(list nodes.List, sep string) string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(nodes.String); ok {
			items = append(items, n.Str)
		}
	}
	return strings.Join(items, sep)
}

func stringSlice(list nodes.List) []string {
	items := []string{}
	for _, item := range list.Items {
		if n, ok := item.(nodes.String); ok {
			items = append(items, n.Str)
		}
	}
	return items
}

type Parameter struct {
	Number int
	Column core.Column
}

// Name and Cmd may be empty
// Maybe I don't need the SQL string if I have the raw Stmt?
type Query struct {
	SQL      string
	Columns  []core.Column
	Params   []Parameter
	Name     string
	Cmd      string // TODO: Pick a better name. One of: one, many, exec, execrows
	Comments []string

	// XXX: Hack
	Filename string
}

type Result struct {
	Queries []*Query
	Catalog core.Catalog
}

type ParserOpts struct {
	UsePositionalParameters bool
}

func ParseQueries(c core.Catalog, queries string, opts ParserOpts) (*Result, error) {
	f, err := os.Stat(queries)
	if err != nil {
		return nil, fmt.Errorf("path %s does not exist", queries)
	}

	var files []string
	if f.IsDir() {
		listing, err := ioutil.ReadDir(queries)
		if err != nil {
			return nil, err
		}
		for _, f := range listing {
			files = append(files, filepath.Join(queries, f.Name()))
		}
	} else {
		files = append(files, queries)
	}

	merr := NewParserErr()
	var q []*Query
	set := map[string]struct{}{}
	for _, filename := range files {
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}
		if strings.HasPrefix(filepath.Base(filename), ".") {
			continue
		}
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		source := string(blob)
		tree, err := pg.Parse(source)
		if err != nil {
			merr.Add(filename, source, 0, err)
			continue
		}
		for _, stmt := range tree.Statements {
			query, err := parseQuery(c, stmt, source, opts.UsePositionalParameters)
			if err == errUnsupportedStatementType {
				continue
			}
			if err != nil {
				merr.Add(filename, source, location(stmt), err)
				continue
			}
			if query.Name != "" {
				if _, exists := set[query.Name]; exists {
					merr.Add(filename, source, location(stmt), fmt.Errorf("duplicate query name: %s", query.Name))
					continue
				}
				set[query.Name] = struct{}{}
			}
			query.Filename = filepath.Base(filename)
			if query != nil {
				q = append(q, query)
			}
		}
	}
	if len(merr.Errs) > 0 {
		return nil, merr
	}
	if len(q) == 0 {
		return nil, fmt.Errorf("path %s contains no queries", queries)
	}
	return &Result{
		Catalog: c,
		Queries: q,
	}, nil
}

func location(node nodes.Node) int {
	switch n := node.(type) {
	case nodes.Query:
		return n.StmtLocation
	case nodes.RawStmt:
		return n.StmtLocation
	}
	return 0
}

func lineno(source string, head int) (int, int) {
	// Calculate the true line and column number for a query, ignoring spaces
	var comment bool
	var loc, line, col int
	for i, char := range source {
		loc += 1
		col += 1
		// TODO: Check bounds
		if char == '-' && source[i+1] == '-' {
			comment = true
		}
		if char == '\n' {
			comment = false
			line += 1
			col = 0
		}
		if loc <= head {
			continue
		}
		if unicode.IsSpace(char) {
			continue
		}
		if comment {
			continue
		}
		break
	}
	return line + 1, col
}

func pluckQuery(source string, n nodes.RawStmt) (string, error) {
	head := n.StmtLocation
	tail := n.StmtLocation + n.StmtLen
	return source[head:tail], nil
}

func rangeVars(root nodes.Node) []nodes.RangeVar {
	var vars []nodes.RangeVar
	find := ast.VisitorFunc(func(node nodes.Node) {
		switch n := node.(type) {
		case nodes.RangeVar:
			vars = append(vars, n)
		}
	})
	ast.Walk(find, root)
	return vars
}

// A query name must be a valid Go identifier
//
// https://golang.org/ref/spec#Identifiers
func validateQueryName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("invalid query name: %q", name)
	}
	for i, c := range name {
		isLetter := unicode.IsLetter(c) || c == '_'
		isDigit := unicode.IsDigit(c)
		if i == 0 && !isLetter {
			return fmt.Errorf("invalid query name %q", name)
		} else if !(isLetter || isDigit) {
			return fmt.Errorf("invalid query name %q", name)
		}
	}
	return nil
}

type CommentSyntax int

const (
	CommentSyntaxDash CommentSyntax = iota
	CommentSyntaxStar               // Note: this is the only style supported by the MySQL sqlparser
	CommentSyntaxHash
)

func ParseMetadata(t string, commentStyle CommentSyntax) (string, string, error) {
	for _, line := range strings.Split(t, "\n") {
		if commentStyle == CommentSyntaxDash && !strings.HasPrefix(line, "-- name:") {
			continue
		}
		if commentStyle == CommentSyntaxStar && !strings.HasPrefix(line, "/* name:") {
			continue
		}
		part := strings.Split(strings.TrimSpace(line), " ")

		if commentStyle == CommentSyntaxStar {
			part = part[:len(part)-1] // removes the trailing "*/" element
		}
		if len(part) == 2 {
			return "", "", fmt.Errorf("missing query type [':one', ':many', ':exec', ':execrows']: %s", line)
		}
		if len(part) != 4 {
			return "", "", fmt.Errorf("invalid query comment: %s", line)
		}
		queryName := part[2]
		queryType := strings.TrimSpace(part[3])
		switch queryType {
		case ":one", ":many", ":exec", ":execrows":
		default:
			return "", "", fmt.Errorf("invalid query type: %s", queryType)
		}
		if err := validateQueryName(queryName); err != nil {
			return "", "", err
		}
		return queryName, queryType, nil
	}
	return "", "", nil
}

func validateCmd(n nodes.Node, name, cmd string) error {
	// TODO: Convert cmd to an enum
	if !(cmd == ":many" || cmd == ":one") {
		return nil
	}
	var list nodes.List
	switch stmt := n.(type) {
	case nodes.SelectStmt:
		return nil
	case nodes.DeleteStmt:
		list = stmt.ReturningList
	case nodes.InsertStmt:
		list = stmt.ReturningList
	case nodes.UpdateStmt:
		list = stmt.ReturningList
	default:
		return nil
	}
	if len(list.Items) == 0 {
		return fmt.Errorf("query %q specifies parameter %q without containing a RETURNING clause", name, cmd)
	}
	return nil
}

var errUnsupportedStatementType = errors.New("parseQuery: unsupported statement type")

func parseQuery(c core.Catalog, stmt nodes.Node, source string, rewriteParameters bool) (*Query, error) {
	if err := validateParamRef(stmt); err != nil {
		return nil, err
	}
	if err := validateParamStyle(stmt); err != nil {
		return nil, err
	}
	raw, ok := stmt.(nodes.RawStmt)
	if !ok {
		return nil, errors.New("node is not a statement")
	}
	switch n := raw.Stmt.(type) {
	case nodes.SelectStmt:
	case nodes.DeleteStmt:
	case nodes.InsertStmt:
		if err := validateInsertStmt(n); err != nil {
			return nil, err
		}
	case nodes.UpdateStmt:
	default:
		return nil, errUnsupportedStatementType
	}

	rawSQL, err := pluckQuery(source, raw)
	if err != nil {
		return nil, err
	}
	if rawSQL == "" {
		return nil, errors.New("missing semicolon at end of file")
	}
	if err := validateFuncCall(&c, raw); err != nil {
		return nil, err
	}
	name, cmd, err := ParseMetadata(strings.TrimSpace(rawSQL), CommentSyntaxDash)
	if err != nil {
		return nil, err
	}
	if err := validateCmd(raw.Stmt, name, cmd); err != nil {
		return nil, err
	}

	// Re-write query AST
	raw, namedParams, edits := rewriteNamedParameters(raw)
	rvs := rangeVars(raw.Stmt)
	refs := findParameters(raw.Stmt)
	if rewriteParameters {
		edits, err = rewriteNumberedParameters(refs, raw, rawSQL)
		if err != nil {
			return nil, err
		}
	} else {
		refs = uniqueParamRefs(refs)
		sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Number < refs[j].ref.Number })
	}
	params, err := resolveCatalogRefs(c, rvs, refs, namedParams)
	if err != nil {
		return nil, err
	}

	qc, err := buildQueryCatalog(c, raw.Stmt)
	if err != nil {
		return nil, err
	}
	cols, err := outputColumns(qc, raw.Stmt)
	if err != nil {
		return nil, err
	}

	expandEdits, err := expand(qc, raw)
	if err != nil {
		return nil, err
	}
	edits = append(edits, expandEdits...)

	expanded, err := editQuery(rawSQL, edits)
	if err != nil {
		return nil, err
	}

	// If the query string was edited, make sure the syntax is valid
	if expanded != rawSQL {
		if _, err := pg.Parse(expanded); err != nil {
			return nil, fmt.Errorf("edited query syntax is invalid: %w", err)
		}
	}

	trimmed, comments, err := stripComments(strings.TrimSpace(expanded))
	if err != nil {
		return nil, err
	}

	return &Query{
		Cmd:      cmd,
		Comments: comments,
		Name:     name,
		Params:   params,
		Columns:  cols,
		SQL:      trimmed,
	}, nil
}

func rewriteNumberedParameters(refs []paramRef, raw nodes.RawStmt, sql string) ([]edit, error) {
	edits := make([]edit, len(refs))
	for i, ref := range refs {
		edits[i] = edit{
			Location: ref.ref.Location - raw.StmtLocation,
			Old:      fmt.Sprintf("$%d", ref.ref.Number),
			New:      "?",
		}
	}
	return edits, nil
}

func stripComments(sql string) (string, []string, error) {
	s := bufio.NewScanner(strings.NewReader(sql))
	var lines, comments []string
	for s.Scan() {
		if strings.HasPrefix(s.Text(), "-- name:") {
			continue
		}
		if strings.HasPrefix(s.Text(), "--") {
			comments = append(comments, strings.TrimPrefix(s.Text(), "--"))
			continue
		}
		lines = append(lines, s.Text())
	}
	return strings.Join(lines, "\n"), comments, s.Err()
}

type edit struct {
	Location int
	Old      string
	New      string
}

func expand(qc *QueryCatalog, raw nodes.RawStmt) ([]edit, error) {
	list := search(raw, func(node nodes.Node) bool {
		switch node.(type) {
		case nodes.DeleteStmt:
		case nodes.InsertStmt:
		case nodes.SelectStmt:
		case nodes.UpdateStmt:
		default:
			return false
		}
		return true
	})
	if len(list.Items) == 0 {
		return nil, nil
	}
	var edits []edit
	for _, item := range list.Items {
		edit, err := expandStmt(qc, raw, item)
		if err != nil {
			return nil, err
		}
		edits = append(edits, edit...)
	}
	return edits, nil
}

func expandStmt(qc *QueryCatalog, raw nodes.RawStmt, node nodes.Node) ([]edit, error) {
	tables, err := sourceTables(qc, node)
	if err != nil {
		return nil, err
	}

	var targets nodes.List
	switch n := node.(type) {
	case nodes.DeleteStmt:
		targets = n.ReturningList
	case nodes.InsertStmt:
		targets = n.ReturningList
	case nodes.SelectStmt:
		targets = n.TargetList
	case nodes.UpdateStmt:
		targets = n.ReturningList
	default:
		return nil, fmt.Errorf("outputColumns: unsupported node type: %T", n)
	}

	var edits []edit
	for _, target := range targets.Items {
		res, ok := target.(nodes.ResTarget)
		if !ok {
			continue
		}
		ref, ok := res.Val.(nodes.ColumnRef)
		if !ok {
			continue
		}
		if !HasStarRef(ref) {
			continue
		}
		var parts, cols []string
		for _, f := range ref.Fields.Items {
			switch field := f.(type) {
			case nodes.String:
				parts = append(parts, field.Str)
			case nodes.A_Star:
				parts = append(parts, "*")
			default:
				return nil, fmt.Errorf("unknown field in ColumnRef: %T", f)
			}
		}
		scope := join(ref.Fields, ".")
		counts := map[string]int{}
		if scope == "" {
			for _, t := range tables {
				for _, c := range t.Columns {
					counts[c.Name] += 1
				}
			}
		}
		for _, t := range tables {
			if scope != "" && scope != t.Name {
				continue
			}
			for _, c := range t.Columns {
				cname := c.Name
				if res.Name != nil {
					cname = *res.Name
				}
				if scope != "" {
					cname = scope + "." + cname
				}
				if counts[cname] > 1 {
					cname = t.Name + "." + cname
				}
				if postgres.IsReservedKeyword(cname) {
					cname = "\"" + cname + "\""
				}
				cols = append(cols, cname)
			}
		}
		edits = append(edits, edit{
			Location: res.Location - raw.StmtLocation,
			Old:      strings.Join(parts, "."),
			New:      strings.Join(cols, ", "),
		})
	}
	return edits, nil
}

func editQuery(raw string, a []edit) (string, error) {
	if len(a) == 0 {
		return raw, nil
	}
	sort.Slice(a, func(i, j int) bool { return a[i].Location > a[j].Location })
	s := raw
	for _, edit := range a {
		start := edit.Location
		if start > len(s) {
			return "", fmt.Errorf("edit start location is out of bounds")
		}
		if len(edit.New) <= 0 {
			return "", fmt.Errorf("empty edit contents")
		}
		if len(edit.Old) <= 0 {
			return "", fmt.Errorf("empty edit contents")
		}
		stop := edit.Location + len(edit.Old) - 1 // Assumes edit.New is non-empty
		if stop < len(s) {
			s = s[:start] + edit.New + s[stop+1:]
		} else {
			s = s[:start] + edit.New
		}
	}
	return s, nil
}

type QueryCatalog struct {
	catalog core.Catalog
	ctes    map[string]core.Table
}

func buildQueryCatalog(c core.Catalog, node nodes.Node) (*QueryCatalog, error) {
	var with *nodes.WithClause
	switch n := node.(type) {
	case nodes.InsertStmt:
		with = n.WithClause
	case nodes.UpdateStmt:
		with = n.WithClause
	case nodes.SelectStmt:
		with = n.WithClause
	default:
		with = nil
	}
	qc := &QueryCatalog{catalog: c, ctes: map[string]core.Table{}}
	if with != nil {
		for _, item := range with.Ctes.Items {
			if cte, ok := item.(nodes.CommonTableExpr); ok {
				cols, err := outputColumns(qc, cte.Ctequery)
				if err != nil {
					return nil, err
				}
				qc.ctes[*cte.Ctename] = core.Table{
					Name:    *cte.Ctename,
					Columns: cols,
				}
			}
		}
	}
	return qc, nil
}

func (qc QueryCatalog) GetTable(fqn core.FQN) (core.Table, *core.Error) {
	cte, exists := qc.ctes[fqn.Rel]
	if exists {
		return cte, nil
	}
	schema, exists := qc.catalog.Schemas[fqn.Schema]
	if !exists {
		err := core.ErrorSchemaDoesNotExist(fqn.Schema)
		return core.Table{}, &err
	}
	table, exists := schema.Tables[fqn.Rel]
	if !exists {
		err := core.ErrorRelationDoesNotExist(fqn.Rel)
		return core.Table{}, &err
	}
	table.ID = fqn
	return table, nil
}

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
// Return an error if a table is referenced twice
// Return an error if an unknown column is referenced
func sourceTables(qc *QueryCatalog, node nodes.Node) ([]core.Table, error) {
	var list nodes.List
	switch n := node.(type) {
	case nodes.DeleteStmt:
		list = nodes.List{
			Items: []nodes.Node{*n.Relation},
		}
	case nodes.InsertStmt:
		list = nodes.List{
			Items: []nodes.Node{*n.Relation},
		}
	case nodes.UpdateStmt:
		list = nodes.List{
			Items: append(n.FromClause.Items, *n.Relation),
		}
	case nodes.SelectStmt:
		list = search(n.FromClause, func(node nodes.Node) bool {
			_, ok := node.(nodes.RangeVar)
			return ok
		})
	default:
		return nil, fmt.Errorf("sourceTables: unsupported node type: %T", n)
	}

	var tables []core.Table
	for _, item := range list.Items {
		switch n := item.(type) {
		case nodes.RangeVar:
			fqn, err := catalog.ParseRange(&n)
			if err != nil {
				return nil, err
			}
			table, cerr := qc.GetTable(fqn)
			if cerr != nil {
				cerr.Location = n.Location
				return nil, *cerr
			}
			if n.Alias != nil {
				table.Name = *n.Alias.Aliasname
			}
			tables = append(tables, table)
		default:
			return nil, fmt.Errorf("sourceTable: unsupported list item type: %T", n)
		}
	}
	return tables, nil
}

func HasStarRef(cf nodes.ColumnRef) bool {
	for _, item := range cf.Fields.Items {
		if _, ok := item.(nodes.A_Star); ok {
			return true
		}
	}
	return false
}

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
func outputColumns(qc *QueryCatalog, node nodes.Node) ([]core.Column, error) {
	tables, err := sourceTables(qc, node)
	if err != nil {
		return nil, err
	}

	var targets nodes.List
	switch n := node.(type) {
	case nodes.DeleteStmt:
		targets = n.ReturningList
	case nodes.InsertStmt:
		targets = n.ReturningList
	case nodes.SelectStmt:
		targets = n.TargetList
	case nodes.UpdateStmt:
		targets = n.ReturningList
	default:
		return nil, fmt.Errorf("outputColumns: unsupported node type: %T", n)
	}

	var cols []core.Column

	for _, target := range targets.Items {
		res, ok := target.(nodes.ResTarget)
		if !ok {
			continue
		}
		switch n := res.Val.(type) {

		case nodes.A_Expr:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			switch {
			case postgres.IsComparisonOperator(join(n.Name, "")):
				// TODO: Generate a name for these operations
				cols = append(cols, core.Column{Name: name, DataType: "bool", NotNull: true})
			case postgres.IsMathematicalOperator(join(n.Name, "")):
				// TODO: Generate correct numeric type
				cols = append(cols, core.Column{Name: name, DataType: "pg_catalog.int4", NotNull: true})
			default:
				cols = append(cols, core.Column{Name: name, DataType: "any", NotNull: false})
			}

		case nodes.CaseExpr:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			// TODO: The TypeCase code has been copied from below. Instead, we need a recurse function to get the type of a node.
			if tc, ok := n.Defresult.(nodes.TypeCast); ok {
				if tc.TypeName == nil {
					return nil, errors.New("no type name type cast")
				}
				name := ""
				if ref, ok := tc.Arg.(nodes.ColumnRef); ok {
					name = join(ref.Fields, "_")
				}
				if res.Name != nil {
					name = *res.Name
				}
				// TODO Validate column names
				col := catalog.ToColumn(tc.TypeName)
				col.Name = name
				cols = append(cols, col)
			} else {
				cols = append(cols, core.Column{Name: name, DataType: "any", NotNull: false})
			}

		case nodes.CoalesceExpr:
			for _, arg := range n.Args.Items {
				if ref, ok := arg.(nodes.ColumnRef); ok {
					columns, err := outputColumnRefs(res, tables, ref)
					if err != nil {
						return nil, err
					}
					for _, c := range columns {
						c.NotNull = true
						cols = append(cols, c)
					}
				}
			}

		case nodes.ColumnRef:
			if HasStarRef(n) {
				// TODO: This code is copied in func expand()
				for _, t := range tables {
					scope := join(n.Fields, ".")
					if scope != "" && scope != t.Name {
						continue
					}
					for _, c := range t.Columns {
						cname := c.Name
						if res.Name != nil {
							cname = *res.Name
						}
						cols = append(cols, core.Column{
							Table:    t.ID,
							Name:     cname,
							Scope:    scope,
							DataType: c.DataType,
							NotNull:  c.NotNull,
							IsArray:  c.IsArray,
						})
					}
				}
				continue
			}

			columns, err := outputColumnRefs(res, tables, n)
			if err != nil {
				return nil, err
			}
			cols = append(cols, columns...)

		case nodes.FuncCall:
			fqn, err := catalog.ParseList(n.Funcname)
			if err != nil {
				return nil, err
			}

			name := fqn.Rel
			if res.Name != nil {
				name = *res.Name
			}

			fun, err := qc.catalog.LookupFunctionN(fqn, len(n.Args.Items))
			if err == nil {
				cols = append(cols, core.Column{Name: name, DataType: fun.ReturnType, NotNull: true})
			} else {
				cols = append(cols, core.Column{Name: name, DataType: "any"})
			}

		case nodes.TypeCast:
			if n.TypeName == nil {
				return nil, errors.New("no type name type cast")
			}
			name := ""
			if ref, ok := n.Arg.(nodes.ColumnRef); ok {
				name = join(ref.Fields, "_")
			}
			if res.Name != nil {
				name = *res.Name
			}
			// TODO Validate column names
			col := catalog.ToColumn(n.TypeName)
			col.Name = name
			cols = append(cols, col)

		default:
			name := ""
			if res.Name != nil {
				name = *res.Name
			}
			cols = append(cols, core.Column{Name: name, DataType: "any", NotNull: false})

		}
	}
	return cols, nil
}

func outputColumnRefs(res nodes.ResTarget, tables []core.Table, node nodes.ColumnRef) ([]core.Column, error) {
	parts := stringSlice(node.Fields)
	var name, alias string
	switch {
	case len(parts) == 1:
		name = parts[0]
	case len(parts) == 2:
		alias = parts[0]
		name = parts[1]
	default:
		return nil, fmt.Errorf("unknown number of fields: %d", len(parts))
	}

	var cols []core.Column
	var found int
	for _, t := range tables {
		if alias != "" && t.Name != alias {
			continue
		}
		for _, c := range t.Columns {
			if c.Name == name {
				found += 1
				cname := c.Name
				if res.Name != nil {
					cname = *res.Name
				}
				cols = append(cols, core.Column{
					Table:    t.ID,
					Name:     cname,
					DataType: c.DataType,
					NotNull:  c.NotNull,
					IsArray:  c.IsArray,
				})
			}
		}
	}
	if found == 0 {
		return nil, core.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column \"%s\" does not exist", name),
			Location: res.Location,
		}
	}
	if found > 1 {
		return nil, core.Error{
			Code:     "42703",
			Message:  fmt.Sprintf("column reference \"%s\" is ambiguous", name),
			Location: res.Location,
		}
	}

	return cols, nil
}

type paramRef struct {
	parent nodes.Node
	rv     *nodes.RangeVar
	ref    nodes.ParamRef
	name   string // Named parameter support
}

type paramSearch struct {
	parent   nodes.Node
	rangeVar *nodes.RangeVar
	refs     *[]paramRef
	seen     map[int]struct{}

	// XXX: Gross state hack for limit
	limitCount  nodes.Node
	limitOffset nodes.Node
}

type nodeImpl struct {
}

func (n nodeImpl) Deparse() string {
	panic("does not deparse")
}

func (n nodeImpl) Fingerprint(nodes.FingerprintContext, nodes.Node, string) {
	panic("does not fingerprint")
}

type limitCount struct {
	nodeImpl
}

type limitOffset struct {
	nodeImpl
}

func (p paramSearch) Visit(node nodes.Node) ast.Visitor {
	switch n := node.(type) {

	case nodes.A_Expr:
		p.parent = node

	case nodes.FuncCall:
		p.parent = node

	case nodes.InsertStmt:
		if s, ok := n.SelectStmt.(nodes.SelectStmt); ok {
			for i, item := range s.TargetList.Items {
				target, ok := item.(nodes.ResTarget)
				if !ok {
					continue
				}
				ref, ok := target.Val.(nodes.ParamRef)
				if !ok {
					continue
				}
				// TODO: Out-of-bounds panic
				*p.refs = append(*p.refs, paramRef{parent: n.Cols.Items[i], ref: ref, rv: n.Relation})
				p.seen[ref.Location] = struct{}{}
			}
			for _, vl := range s.ValuesLists {
				for i, v := range vl {
					ref, ok := v.(nodes.ParamRef)
					if !ok {
						continue
					}
					// TODO: Out-of-bounds panic
					*p.refs = append(*p.refs, paramRef{parent: n.Cols.Items[i], ref: ref, rv: n.Relation})
					p.seen[ref.Location] = struct{}{}
				}
			}
		}

	case nodes.RangeVar:
		p.rangeVar = &n

	case nodes.ResTarget:
		p.parent = node

	case nodes.SelectStmt:
		if n.LimitCount != nil {
			p.limitCount = n.LimitCount
		}
		if n.LimitOffset != nil {
			p.limitOffset = n.LimitOffset
		}

	case nodes.TypeCast:
		p.parent = node

	case nodes.ParamRef:
		parent := p.parent

		if count, ok := p.limitCount.(nodes.ParamRef); ok {
			if n.Number == count.Number {
				parent = limitCount{}
			}
		}

		if offset, ok := p.limitOffset.(nodes.ParamRef); ok {
			if n.Number == offset.Number {
				parent = limitOffset{}
			}
		}
		if _, found := p.seen[n.Location]; found {
			break
		}

		// Special, terrible case for nodes.MultiAssignRef
		set := true
		if res, ok := parent.(nodes.ResTarget); ok {
			if multi, ok := res.Val.(nodes.MultiAssignRef); ok {
				set = false
				if row, ok := multi.Source.(nodes.RowExpr); ok {
					for i, arg := range row.Args.Items {
						if ref, ok := arg.(nodes.ParamRef); ok {
							if multi.Colno == i+1 && ref.Number == n.Number {
								set = true
							}
						}
					}
				}
			}
		}

		if set {
			*p.refs = append(*p.refs, paramRef{parent: parent, ref: n, rv: p.rangeVar})
			p.seen[n.Location] = struct{}{}
		}
		return nil
	}
	return p
}

func findParameters(root nodes.Node) []paramRef {
	refs := make([]paramRef, 0)
	v := paramSearch{seen: make(map[int]struct{}), refs: &refs}
	ast.Walk(v, root)
	return refs
}

type nodeSearch struct {
	list  nodes.List
	check func(nodes.Node) bool
}

func (s *nodeSearch) Visit(node nodes.Node) ast.Visitor {
	if s.check(node) {
		s.list.Items = append(s.list.Items, node)
	}
	return s
}

func search(root nodes.Node, f func(nodes.Node) bool) nodes.List {
	ns := &nodeSearch{check: f}
	ast.Walk(ns, root)
	return ns.list
}

func resolveCatalogRefs(c core.Catalog, rvs []nodes.RangeVar, args []paramRef, names map[int]string) ([]Parameter, error) {
	aliasMap := map[string]core.FQN{}
	// TODO: Deprecate defaultTable
	var defaultTable *core.FQN
	var tables []core.FQN

	parameterName := func(n int, defaultName string) string {
		if n, ok := names[n]; ok {
			return n
		}
		return defaultName
	}

	for _, rv := range rvs {
		if rv.Relname == nil {
			continue
		}
		fqn, err := catalog.ParseRange(&rv)
		if err != nil {
			return nil, err
		}
		tables = append(tables, fqn)
		if defaultTable == nil {
			defaultTable = &fqn
		}
		if rv.Alias == nil {
			continue
		}
		aliasMap[*rv.Alias.Aliasname] = fqn
	}

	typeMap := map[string]map[string]map[string]core.Column{}
	for _, fqn := range tables {
		schema, found := c.Schemas[fqn.Schema]
		if !found {
			continue
		}

		table, found := schema.Tables[fqn.Rel]
		if !found {
			continue
		}

		if _, exists := typeMap[fqn.Schema]; !exists {
			typeMap[fqn.Schema] = map[string]map[string]core.Column{}
		}

		typeMap[fqn.Schema][fqn.Rel] = map[string]core.Column{}
		for _, c := range table.Columns {
			cc := c
			typeMap[fqn.Schema][fqn.Rel][c.Name] = cc
		}
	}

	var a []Parameter
	for _, ref := range args {
		switch n := ref.parent.(type) {

		case limitOffset:
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: core.Column{
					Name:     parameterName(ref.ref.Number, "offset"),
					DataType: "integer",
					NotNull:  true,
				},
			})

		case limitCount:
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: core.Column{
					Name:     parameterName(ref.ref.Number, "limit"),
					DataType: "integer",
					NotNull:  true,
				},
			})

		case nodes.A_Expr:
			// TODO: While this works for a wide range of simple expressions,
			// more complicated expressions will cause this logic to fail.
			list := search(n.Lexpr, func(node nodes.Node) bool {
				_, ok := node.(nodes.ColumnRef)
				return ok
			})

			if len(list.Items) == 0 {
				return nil, core.Error{
					Code:     "XXXXX",
					Message:  "no column reference found",
					Location: n.Location,
				}
			}

			switch left := list.Items[0].(type) {
			case nodes.ColumnRef:
				items := stringSlice(left.Fields)
				var key, alias string
				switch len(items) {
				case 1:
					key = items[0]
				case 2:
					alias = items[0]
					key = items[1]
				default:
					panic("too many field items: " + strconv.Itoa(len(items)))
				}

				search := tables
				if alias != "" {
					if original, ok := aliasMap[alias]; ok {
						search = []core.FQN{original}
					} else {
						for _, fqn := range tables {
							if fqn.Rel == alias {
								search = []core.FQN{fqn}
							}
						}
					}
				}

				var found int
				for _, table := range search {
					if c, ok := typeMap[table.Schema][table.Rel][key]; ok {
						found += 1
						if ref.name != "" {
							key = ref.name
						}
						a = append(a, Parameter{
							Number: ref.ref.Number,
							Column: core.Column{
								Name:     parameterName(ref.ref.Number, key),
								DataType: c.DataType,
								NotNull:  c.NotNull,
								IsArray:  c.IsArray,
								Table:    c.Table,
							},
						})
					}
				}
				if found == 0 {
					return nil, core.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column \"%s\" does not exist", key),
						Location: left.Location,
					}
				}
				if found > 1 {
					return nil, core.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column reference \"%s\" is ambiguous", key),
						Location: left.Location,
					}
				}
			}

		case nodes.FuncCall:
			fqn, err := catalog.ParseList(n.Funcname)
			if err != nil {
				return nil, err
			}
			fun, err := c.LookupFunctionN(fqn, len(n.Args.Items))
			if err != nil {
				// Synthesize a function on the fly to avoid returning with an error
				// for an unknown Postgres function (e.g. defined in an extension)
				fun = core.Function{
					Name:       fqn.Rel,
					ArgN:       len(n.Args.Items),
					Arguments:  nil,
					ReturnType: "any",
				}
			}
			for i, item := range n.Args.Items {
				pr, ok := item.(nodes.ParamRef)
				if !ok {
					continue
				}
				if pr.Number != ref.ref.Number {
					continue
				}
				if fun.Arguments == nil {
					a = append(a, Parameter{
						Number: ref.ref.Number,
						Column: core.Column{
							Name:     parameterName(ref.ref.Number, fun.Name),
							DataType: "any",
						},
					})
					continue
				}
				if i >= len(fun.Arguments) {
					return nil, fmt.Errorf("incorrect number of arguments to %s", fun.Name)
				}
				arg := fun.Arguments[i]
				name := arg.Name
				if name == "" {
					name = fun.Name
				}
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: core.Column{
						Name:     parameterName(ref.ref.Number, name),
						DataType: arg.DataType,
						NotNull:  true,
					},
				})
			}

		case nodes.ResTarget:
			if n.Name == nil {
				return nil, fmt.Errorf("nodes.ResTarget has nil name")
			}
			key := *n.Name

			// TODO: Deprecate defaultTable
			schema := defaultTable.Schema
			rel := defaultTable.Rel
			if ref.rv != nil {
				fqn, err := catalog.ParseRange(ref.rv)
				if err != nil {
					return nil, err
				}
				schema = fqn.Schema
				rel = fqn.Rel
			}
			if c, ok := typeMap[schema][rel][key]; ok {
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: core.Column{
						Name:     parameterName(ref.ref.Number, key),
						DataType: c.DataType,
						NotNull:  c.NotNull,
						IsArray:  c.IsArray,
						Table:    c.Table,
					},
				})
			} else {
				return nil, core.Error{
					Code:     "42703",
					Message:  fmt.Sprintf("column \"%s\" does not exist", key),
					Location: n.Location,
				}
			}

		case nodes.TypeCast:
			if n.TypeName == nil {
				return nil, fmt.Errorf("nodes.TypeCast has nil type name")
			}
			col := catalog.ToColumn(n.TypeName)
			col.Name = parameterName(ref.ref.Number, col.Name)
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: col,
			})

		case nodes.ParamRef:
			a = append(a, Parameter{Number: ref.ref.Number})

		default:
			fmt.Printf("unsupported reference type: %T", n)
		}
	}
	return a, nil
}

func uniqueParamRefs(in []paramRef) []paramRef {
	m := make(map[int]struct{}, len(in))
	o := make([]paramRef, 0, len(in))
	for _, v := range in {
		if _, ok := m[v.ref.Number]; !ok {
			m[v.ref.Number] = struct{}{}
			o = append(o, v)
		}
	}
	return o
}
