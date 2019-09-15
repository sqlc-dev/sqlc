package dinosql

import (
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

func ParseCatalog(schema string) (core.Catalog, error) {
	f, err := os.Stat(schema)
	if err != nil {
		return core.Catalog{}, fmt.Errorf("path %s does not exist", schema)
	}

	var files []os.FileInfo
	if f.IsDir() {
		files, err = ioutil.ReadDir(schema)
		if err != nil {
			return core.Catalog{}, err
		}
	} else {
		files = append(files, f)
	}

	merr := NewParserErr()
	c := core.NewCatalog()
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		filename := filepath.Join(schema, f.Name())
		blob, err := ioutil.ReadFile(filename)
		if err != nil {
			merr.Add(filename, "", 0, err)
			continue
		}
		contents := RemoveGooseRollback(string(blob))
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
	SQL     string
	Columns []core.Column
	Params  []Parameter
	Name    string
	Cmd     string // TODO: Pick a better name. One of: one, many, exec, execrows

	// XXX: Hack
	NeedsEdit bool
	Filename  string
}

type Result struct {
	Settings GenerateSettings
	Queries  []*Query
	Catalog  core.Catalog
}

func ParseQueries(c core.Catalog, settings GenerateSettings, pkg PackageSettings) (*Result, error) {
	f, err := os.Stat(pkg.Queries)
	if err != nil {
		return nil, fmt.Errorf("path %s does not exist", pkg.Queries)
	}

	var files []os.FileInfo
	if f.IsDir() {
		files, err = ioutil.ReadDir(pkg.Queries)
		if err != nil {
			return nil, err
		}
	} else {
		files = append(files, f)
	}

	merr := NewParserErr()
	var q []*Query
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		filename := filepath.Join(pkg.Queries, f.Name())
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
			// line, col := location(source, stmt)
			query, err := parseQuery(c, stmt, source)
			if err != nil {
				merr.Add(filename, source, location(stmt), err)
				continue
			}
			query.Filename = f.Name()
			if query != nil {
				q = append(q, query)
			}
		}
	}

	if len(merr.Errs) > 0 {
		return nil, merr
	}

	return &Result{Catalog: c, Queries: q, Settings: settings}, nil
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
	return strings.TrimSpace(source[head:tail]), nil
}

func rangeVars(root nodes.Node) []nodes.RangeVar {
	var vars []nodes.RangeVar
	find := VisitorFunc(func(node nodes.Node) {
		switch n := node.(type) {
		case nodes.RangeVar:
			vars = append(vars, n)
		}
	})
	Walk(find, root)
	return vars
}

// TODO: Validate metadata
func parseMetadata(t string) (string, string, error) {
	for _, line := range strings.Split(t, "\n") {
		if !strings.HasPrefix(line, "-- name:") {
			continue
		}
		part := strings.Split(line, " ")
		return part[2], strings.TrimSpace(part[3]), nil
	}
	return "", "", nil
}

func parseQuery(c core.Catalog, stmt nodes.Node, source string) (*Query, error) {
	if err := validateParamRef(stmt); err != nil {
		return nil, err
	}
	raw, ok := stmt.(nodes.RawStmt)
	if !ok {
		return nil, nil
	}
	switch raw.Stmt.(type) {
	case nodes.SelectStmt:
	case nodes.DeleteStmt:
	case nodes.InsertStmt:
	case nodes.UpdateStmt:
	default:
		return nil, nil
	}
	rawSQL, err := pluckQuery(source, raw)
	if err != nil {
		return nil, err
	}
	if err := validateFuncCall(&c, raw); err != nil {
		return nil, err
	}
	name, cmd, err := parseMetadata(rawSQL)
	if err != nil {
		return nil, err
	}

	rvs := rangeVars(raw.Stmt)
	refs := findParameters(raw.Stmt)
	params, err := resolveCatalogRefs(c, rvs, refs)
	if err != nil {
		return nil, err
	}

	cols, err := outputColumns(c, raw.Stmt)
	if err != nil {
		return nil, err
	}

	return &Query{
		Cmd:       cmd,
		Name:      name,
		Params:    params,
		Columns:   cols,
		SQL:       rawSQL,
		NeedsEdit: needsEdit(stmt),
	}, nil
}

type QueryCatalog struct {
	catalog core.Catalog
	ctes    map[string]core.Table
}

func NewQueryCatalog(c core.Catalog, with *nodes.WithClause) QueryCatalog {
	ctes := map[string]core.Table{}
	if with != nil {
		for _, item := range with.Ctes.Items {
			if cte, ok := item.(nodes.CommonTableExpr); ok {
				cols, err := outputColumns(c, cte.Ctequery)
				if err != nil {
					panic(err.Error())
				}
				ctes[*cte.Ctename] = core.Table{
					Name:    *cte.Ctename,
					Columns: cols,
				}
			}
		}
	}
	return QueryCatalog{catalog: c, ctes: ctes}
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
	return table, nil
}

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
// Return an error if a table is referenced twice
// Return an error if an unknown column is referenced
func sourceTables(c core.Catalog, node nodes.Node) ([]core.Table, error) {
	var list nodes.List
	var with *nodes.WithClause
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
		with = n.WithClause
		list = search(n.FromClause, func(node nodes.Node) bool {
			_, ok := node.(nodes.RangeVar)
			return ok
		})
	default:
		return nil, fmt.Errorf("sourceTables: unsupported node type: %T", n)
	}

	qc := NewQueryCatalog(c, with)

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
func outputColumns(c core.Catalog, node nodes.Node) ([]core.Column, error) {
	tables, err := sourceTables(c, node)
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
		// spew.Dump(target)

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
			if postgres.IsComparisonOperator(join(n.Name, "")) {
				// TODO: Generate a name for these operations
				cols = append(cols, core.Column{Name: name, DataType: "bool", NotNull: true})
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

			fun, err := c.LookupFunctionN(fqn, len(n.Args.Items))
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
					Name:     cname,
					DataType: c.DataType,
					NotNull:  c.NotNull,
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
}

type paramSearch struct {
	parent   nodes.Node
	rangeVar *nodes.RangeVar
	refs     map[int]paramRef

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

func (p *paramSearch) Visit(node nodes.Node) Visitor {
	switch n := node.(type) {

	case nodes.A_Expr:
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
				p.refs[ref.Number] = paramRef{parent: n.Cols.Items[i], ref: ref, rv: p.rangeVar}
			}
			for _, vl := range s.ValuesLists {
				for i, v := range vl {
					ref, ok := v.(nodes.ParamRef)
					if !ok {
						continue
					}
					// TODO: Out-of-bounds panic
					p.refs[ref.Number] = paramRef{parent: n.Cols.Items[i], ref: ref, rv: p.rangeVar}
				}
			}
		}

	case nodes.RangeVar:
		p.rangeVar = &n

	case nodes.ResTarget:
		p.parent = node

	case nodes.SelectStmt:
		p.limitCount = n.LimitCount
		p.limitOffset = n.LimitOffset

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
		if _, found := p.refs[n.Number]; found {
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
			p.refs[n.Number] = paramRef{parent: parent, ref: n, rv: p.rangeVar}
		}
		return nil
	}
	return p
}

func findParameters(root nodes.Node) []paramRef {
	v := &paramSearch{refs: map[int]paramRef{}}
	Walk(v, root)
	refs := make([]paramRef, 0)
	for _, r := range v.refs {
		refs = append(refs, r)
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Number < refs[j].ref.Number })
	return refs
}

type starWalker struct {
	found bool
}

func (s *starWalker) Visit(node nodes.Node) Visitor {
	if _, ok := node.(nodes.A_Star); ok {
		s.found = true
		return nil
	}
	return s
}

func needsEdit(root nodes.Node) bool {
	v := &starWalker{}
	Walk(v, root)
	return v.found
}

type nodeSearch struct {
	list  nodes.List
	check func(nodes.Node) bool
}

func (s *nodeSearch) Visit(node nodes.Node) Visitor {
	if s.check(node) {
		s.list.Items = append(s.list.Items, node)
	}
	return s
}

func search(root nodes.Node, f func(nodes.Node) bool) nodes.List {
	ns := &nodeSearch{check: f}
	Walk(ns, root)
	return ns.list
}

func argName(name string) string {
	out := ""
	for i, p := range strings.Split(name, "_") {
		if i == 0 {
			out += strings.ToLower(p)
		} else if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

func resolveCatalogRefs(c core.Catalog, rvs []nodes.RangeVar, args []paramRef) ([]Parameter, error) {
	aliasMap := map[string]core.FQN{}
	// TODO: Deprecate defaultTable
	var defaultTable *core.FQN
	var tables []core.FQN

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
					Name:     "offset",
					DataType: "integer",
					NotNull:  true,
				},
			})

		case limitCount:
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: core.Column{
					Name:     "limit",
					DataType: "integer",
					NotNull:  true,
				},
			})

		case nodes.A_Expr:
			switch n := n.Lexpr.(type) {
			case nodes.ColumnRef:
				items := stringSlice(n.Fields)
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
						a = append(a, Parameter{
							Number: ref.ref.Number,
							Column: core.Column{
								Name:     argName(key),
								DataType: c.DataType,
								NotNull:  c.NotNull,
								IsArray:  c.IsArray,
							},
						})
					}
				}
				if found == 0 {
					return nil, core.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column \"%s\" does not exist", key),
						Location: n.Location,
					}
				}
				if found > 1 {
					return nil, core.Error{
						Code:     "42703",
						Message:  fmt.Sprintf("column reference \"%s\" is ambiguous", key),
						Location: n.Location,
					}
				}
			}

		case nodes.ResTarget:
			if n.Name == nil {
				return nil, fmt.Errorf("nodes.ResTarget has nil name")
			}
			key := *n.Name
			if c, ok := typeMap[defaultTable.Schema][defaultTable.Rel][key]; ok {
				a = append(a, Parameter{
					Number: ref.ref.Number,
					Column: core.Column{
						Name:     argName(key),
						DataType: c.DataType,
						NotNull:  c.NotNull,
						IsArray:  c.IsArray,
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
			a = append(a, Parameter{
				Number: ref.ref.Number,
				Column: catalog.ToColumn(n.TypeName),
			})

		case nodes.ParamRef:
			a = append(a, Parameter{Number: ref.ref.Number})

		default:
			// return nil, fmt.Errorf("unsupported type: %T", n)
		}
	}
	return a, nil
}

type TypeOverride struct {
	Package      string `json:"package"`
	PostgresType string `json:"postgres_type"`
	GoType       string `json:"go_type"`
	Null         bool   `json:"null"`
}

type PackageSettings struct {
	Name                string `json:"name"`
	Path                string `json:"path"`
	Schema              string `json:"schema"`
	Queries             string `json:"queries"`
	EmitPreparedQueries bool   `json:"emit_prepared_queries"`
	EmitTags            bool   `json:"emit_tags"`
}

type GenerateSettings struct {
	Packages  []PackageSettings `json:"packages"`
	Overrides []TypeOverride    `json:"overrides"`
	Rename    map[string]string `json:"rename"`
}
