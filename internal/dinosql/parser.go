package dinosql

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/dinosql/internal/catalog"
	core "github.com/kyleconroy/dinosql/internal/pg"
	"github.com/kyleconroy/dinosql/internal/postgres"

	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func keepSpew() {
	spew.Dump("hello world")
}

func ParseCatalog(dir string, settings GenerateSettings) (core.Catalog, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return core.Catalog{}, err
	}
	c := core.NewCatalog()
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		blob, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return c, err
		}
		contents := RemoveGooseRollback(string(blob))
		tree, err := pg.Parse(contents)
		if err != nil {
			return c, err
		}
		if err := updateCatalog(&c, tree); err != nil {
			return c, err
		}
	}
	return c, nil
}

func updateCatalog(c *core.Catalog, tree pg.ParsetreeList) error {
	for _, stmt := range tree.Statements {
		if err := validateFuncCall(stmt); err != nil {
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
}

type Result struct {
	Settings GenerateSettings
	Queries  []*Query
	Catalog  core.Catalog
}

func ParseQueries(c core.Catalog, settings GenerateSettings) (*Result, error) {
	files, err := ioutil.ReadDir(settings.QueryDir)
	if err != nil {
		return nil, err
	}

	var q []*Query
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		blob, err := ioutil.ReadFile(filepath.Join(settings.QueryDir, f.Name()))
		if err != nil {
			return nil, err
		}
		source := string(blob)
		tree, err := pg.Parse(source)
		if err != nil {
			return nil, err
		}
		for _, stmt := range tree.Statements {
			queryTwo, err := parseQuery(c, stmt, source)
			if err != nil {
				return nil, err
			}
			if queryTwo != nil {
				q = append(q, queryTwo)
			}
		}
	}
	return &Result{Catalog: c, Queries: q, Settings: settings}, nil
}

func pluckQuery(source string, n nodes.RawStmt) (string, error) {
	// TODO: Bounds checking
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
	if err := validateFuncCall(raw); err != nil {
		return nil, err
	}
	rawSQL, err := pluckQuery(source, raw)
	if err != nil {
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

func (qc QueryCatalog) GetTable(fqn core.FQN) (core.Table, error) {
	cte, exists := qc.ctes[fqn.Rel]
	if exists {
		return cte, nil
	}
	schema, exists := qc.catalog.Schemas[fqn.Schema]
	if !exists {
		return core.Table{}, core.ErrorSchemaDoesNotExist(fqn.Schema)
	}
	table, exists := schema.Tables[fqn.Rel]
	if !exists {
		return core.Table{}, core.ErrorRelationDoesNotExist(fqn.Rel)
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
		list = n.FromClause
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
			table, err := qc.GetTable(fqn)
			if err != nil {
				return nil, err
			}
			tables = append(tables, table)
		default:
			return nil, fmt.Errorf("sourceTable: unsupported list item type: %T", n)
		}
	}
	return tables, nil
}

func IsStarRef(cf nodes.ColumnRef) bool {
	if len(cf.Fields.Items) != 1 {
		return false
	}
	_, aStar := cf.Fields.Items[0].(nodes.A_Star)
	return aStar
}

// Compute the output columns for a statement.
//
// Return an error if column references are ambiguous
// Return an error if column references don't exist
func outputColumns(c core.Catalog, node nodes.Node) ([]core.Column, error) {
	tables, err := sourceTables(c, node)
	if err != nil {
		fmt.Println(tables)
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
			name := "_"
			if res.Name != nil {
				name = *res.Name
			}
			if postgres.IsComparisonOperator(join(n.Name, "")) {
				// TODO: Generate a name for these operations
				cols = append(cols, core.Column{Name: name, DataType: "bool", NotNull: true})
			}

		case nodes.ColumnRef:
			parts := stringSlice(n.Fields)
			var name, alias string
			switch {
			case IsStarRef(n):
				// TODO: Disambiguate columns
				for _, t := range tables {
					for _, c := range t.Columns {
						cname := c.Name
						if res.Name != nil {
							cname = *res.Name
						}
						cols = append(cols, core.Column{
							Name:     cname,
							DataType: c.DataType,
							NotNull:  c.NotNull,
							IsArray:  c.IsArray,
						})
					}
				}
				continue
			case len(parts) == 1:
				name = parts[0]
			case len(parts) == 2:
				alias = parts[0]
				name = parts[1]
			default:
				panic(fmt.Sprintf("unknown number of fields: %d", len(parts)))
			}
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
				return nil, Error{
					Code:    "42703",
					Message: fmt.Sprintf("column \"%s\" does not exist", name),
				}
			}
			if found > 1 {
				return nil, Error{
					Code:    "42703",
					Message: fmt.Sprintf("column reference \"%s\" is ambiguous", name),
				}
			}

		case nodes.FuncCall:
			// TODO: Look up return type of functions
			name := join(n.Funcname, ".")
			if res.Name != nil {
				name = *res.Name
			}
			cols = append(cols, core.Column{Name: name, DataType: "bigint"})

		case nodes.TypeCast:
			if n.TypeName == nil {
				return nil, errors.New("no type name type cast")
			}
			name := ""
			if ref, ok := n.Arg.(nodes.ColumnRef); ok {
				name = join(ref.Fields, "_")
			}
			// TODO Validate column names
			col := catalog.ToColumn(n.TypeName)
			col.Name = name
			cols = append(cols, col)
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

		if _, found := p.refs[n.Number]; !found {
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
	typeMap := map[string]map[string]core.Column{}
	for _, t := range c.Schemas["public"].Tables {
		typeMap[t.Name] = map[string]core.Column{}
		for _, c := range t.Columns {
			cc := c
			typeMap[t.Name][c.Name] = cc
		}
	}

	aliasMap := map[string]string{}
	defaultTable := ""
	for _, rv := range rvs {
		if rv.Relname == nil {
			continue
		}
		if defaultTable == "" {
			defaultTable = *rv.Relname
		}
		if rv.Alias == nil {
			continue
		}
		aliasMap[*rv.Alias.Aliasname] = *rv.Relname
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

				table := aliasMap[alias]
				if table == "" && ref.rv != nil && ref.rv.Relname != nil {
					table = *ref.rv.Relname
				}
				if table == "" {
					table = defaultTable
				}

				if c, ok := typeMap[table][key]; ok {
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
					return nil, Error{
						Code:    "42703",
						Message: fmt.Sprintf("column \"%s\" does not exist", key),
					}
				}
			}
		case nodes.ResTarget:
			if n.Name == nil {
				return nil, fmt.Errorf("nodes.ResTarget has nil name")
			}
			key := *n.Name
			if c, ok := typeMap[defaultTable][key]; ok {
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
				return nil, Error{
					Code:    "42703",
					Message: fmt.Sprintf("column \"%s\" does not exist", key),
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

type GenerateSettings struct {
	SchemaDir           string            `json:"schema"`
	QueryDir            string            `json:"queries"`
	Out                 string            `json:"out"`
	Package             string            `json:"package"`
	EmitPreparedQueries bool              `json:"emit_prepared_queries"`
	EmitTags            bool              `json:"emit_tags"`
	Overrides           []TypeOverride    `json:"overrides"`
	Rename              map[string]string `json:"rename"`
}
