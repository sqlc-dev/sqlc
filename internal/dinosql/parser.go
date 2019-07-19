package dinosql

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/kyleconroy/dinosql/internal/catalog"
	core "github.com/kyleconroy/dinosql/internal/pg"
	"github.com/kyleconroy/dinosql/internal/postgres"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/inflection"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func parseSQL(in string) (*Result, error) {
	if false {
		spew.Dump(in)
	}
	tree, err := pg.Parse(in)
	if err != nil {
		return nil, err
	}
	c := core.NewCatalog()
	if err := updateCatalog(&c, tree); err != nil {
		return nil, err
	}

	var q []Query
	s := convert(c, GenerateSettings{})
	r := Result{Schema: s}
	if err := parseFuncs(s, &r, in, tree); err != nil {
		return nil, err
	}
	q = append(q, r.Queries...)

	return &Result{Schema: s, Queries: q}, nil
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

// Eventually get rid of the postgres package. But for now, generate a
// postgres.Schema from a pg.Catalog
func convert(c core.Catalog, settings GenerateSettings) *postgres.Schema {
	s := postgres.Schema{}

	for name, schema := range c.Schemas {
		// For now, only convert the public schema
		if name != "public" {
			continue
		}
		for _, enum := range schema.Enums {
			s.Enums = append(s.Enums, postgres.Enum{
				Name:   enum.Name,
				GoName: structName(enum.Name),
				Vals:   enum.Vals,
			})
		}
		for _, table := range schema.Tables {
			t := postgres.Table{
				Name:   table.Name,
				GoName: inflection.Singular(structName(table.Name)),
			}
			for _, column := range table.Columns {
				t.Columns = append(t.Columns, postgres.Column{
					Name:    column.Name,
					Type:    column.DataType,
					NotNull: column.NotNull,
					GoName:  structName(column.Name),
					GoType:  columnType(&s, settings, column.DataType, column.NotNull),
				})
			}
			s.Tables = append(s.Tables, t)
		}
	}

	sort.Slice(s.Tables, func(i, j int) bool { return s.Tables[i].Name < s.Tables[j].Name })
	sort.Slice(s.Enums, func(i, j int) bool { return s.Enums[i].Name < s.Enums[j].Name })
	return &s
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

func isStar(n outputRef) bool {
	if n.ref == nil {
		return false
	}
	if len(n.ref.Fields.Items) != 1 {
		return false
	}
	_, aStar := n.ref.Fields.Items[0].(nodes.A_Star)
	return aStar
}

type Field struct {
	Name string
	Type string
}

type Column struct {
	Name string
	Type string
}

type Parameter struct {
	Number   int
	DataType string
	Name     string // TODO: Relation?
	NotNull  bool
}

// Name and Cmd may be empty
// Maybe I don't need the SQL string if I have the raw Stmt?
type QueryTwo struct {
	SQL     string
	Columns []core.Column
	Outs    []outputRef
	Params  []Parameter
	Name    string
	Cmd     string // TODO: Pick a better name. One of: one, many, exec, execrows
}

// TODO: The Query struct is overloaded
type Query struct {
	Type       string
	MethodName string
	StmtName   string
	QueryName  string
	SQL        string
	Args       []Arg
	Table      postgres.Table
	Fields     []Field
	ReturnType string
	RowStruct  bool
	ScanRecord bool
}

type GoConstant struct {
	Name  string
	Type  string
	Value string
}

type GoEnum struct {
	Name      string
	Constants []GoConstant
}

type GoField struct {
	Name string
	Type string
	Tags map[string]string
}

func (gf GoField) Tag() string {
	var tags []string
	for key, val := range gf.Tags {
		tags = append(tags, fmt.Sprintf("%s\"%s\"", key, val))
	}
	if len(tags) == 0 {
		return ""
	}
	sort.Strings(tags)
	return strings.Join(tags, ",")
}

type GoStruct struct {
	Name   string
	Fields []GoField
}

// TODO: Terrible name
type GoQueryValue struct {
	Emit   bool
	Name   string
	Struct *GoStruct
	typ    string
}

func (v GoQueryValue) EmitStruct() bool {
	return v.Emit
}

func (v GoQueryValue) IsStruct() bool {
	return v.Struct != nil
}

func (v GoQueryValue) isEmpty() bool {
	return v.typ == "" && v.Name == "" && v.Struct == nil
}

func (v GoQueryValue) Pair() string {
	if v.isEmpty() {
		return ""
	}
	return v.Name + " " + v.Type()
}

func (v GoQueryValue) Type() string {
	if v.typ != "" {
		return v.typ
	}
	if v.Struct != nil {
		return v.Struct.Name
	}
	panic("no type for GoQueryValue: " + v.Name)
}

func (v GoQueryValue) Params() string {
	if v.isEmpty() {
		return ""
	}
	var out []string
	if v.Struct == nil {
		out = append(out, v.Name)
	} else {
		for _, f := range v.Struct.Fields {
			out = append(out, v.Name+"."+f.Name)
		}
	}
	return strings.Join(out, ",")
}

func (v GoQueryValue) Scan() string {
	var out []string
	if v.Struct == nil {
		out = append(out, "&"+v.Name)
	} else {
		for _, f := range v.Struct.Fields {
			out = append(out, "&"+v.Name+"."+f.Name)
		}
	}
	return strings.Join(out, ",")
}

// A struct used to generate methods and fields on the Queries struct
type GoQuery struct {
	Cmd          string
	MethodName   string
	FieldName    string
	ConstantName string
	SQL          string
	Ret          GoQueryValue
	Arg          GoQueryValue
}

type Result struct {
	Settings  GenerateSettings
	Schema    *postgres.Schema
	Queries   []Query
	QueryTwos []*QueryTwo
	Catalog   core.Catalog
}

func (r Result) UsesType(typ string) bool {
	for _, strct := range r.Structs() {
		for _, f := range strct.Fields {
			if f.Type == typ {
				return true
			}
		}
	}
	return false
}

func (r Result) StdImports() []string {
	imports := []string{
		"context",
		"database/sql",
	}
	if r.UsesType("json.RawMessage") {
		imports = append(imports, "encoding/json")
	}
	if r.UsesType("time.Time") {
		imports = append(imports, "time")
	}
	return imports
}

func (r Result) PkgImports(settings GenerateSettings) []string {
	imports := []string{}

	if r.UsesType("pq.NullTime") {
		imports = append(imports, "github.com/lib/pq")
	}
	if r.UsesType("uuid.UUID") {
		imports = append(imports, "github.com/google/uuid")
	}

	// Custom imports
	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		overrideTypes[o.GoType] = o.Package
	}
	for goType, importPath := range overrideTypes {
		if r.UsesType(goType) {
			imports = append(imports, importPath)
		}
	}
	return imports
}

func (r Result) Enums() []GoEnum {
	var enums []GoEnum
	for name, schema := range r.Catalog.Schemas {
		if name != "public" {
			continue
		}
		for _, enum := range schema.Enums {
			e := GoEnum{
				Name: structName(enum.Name),
			}
			for _, v := range enum.Vals {
				name := ""
				for _, part := range strings.Split(strings.Replace(v, "-", "_", -1), "_") {
					name += strings.Title(part)
				}
				e.Constants = append(e.Constants, GoConstant{
					Name:  e.Name + name,
					Value: v,
					Type:  e.Name,
				})
			}
			enums = append(enums, e)
		}
	}
	if len(enums) > 0 {
		sort.Slice(enums, func(i, j int) bool { return enums[i].Name < enums[j].Name })
	}
	return enums
}

func (r Result) Structs() []GoStruct {
	var structs []GoStruct
	for name, schema := range r.Catalog.Schemas {
		if name != "public" {
			continue
		}
		for _, table := range schema.Tables {
			s := GoStruct{
				Name: inflection.Singular(structName(table.Name)),
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, GoField{
					Name: structName(column.Name),
					Type: r.goType(column.DataType, column.NotNull),
					Tags: map[string]string{"json": column.Name},
				})
			}
			structs = append(structs, s)
		}
	}
	if len(structs) > 0 {
		sort.Slice(structs, func(i, j int) bool { return structs[i].Name < structs[j].Name })
	}
	return structs
}

func (r Result) goType(columnType string, notNull bool) string {
	for _, oride := range r.Settings.Overrides {
		if oride.PostgresType == columnType && oride.Null != notNull {
			return oride.GoType
		}
	}

	switch columnType {
	case "text":
		if notNull {
			return "string"
		} else {
			return "sql.NullString"
		}
	case "serial":
		return "int"
	case "integer":
		return "int"
	case "bool":
		if notNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}
	case "jsonb":
		return "json.RawMessage"
	case "pg_catalog.bool":
		if notNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}
	case "pg_catalog.int2":
		return "uint8"
	case "pg_catalog.int4":
		return "int"
	case "pg_catalog.int8":
		return "int"
	case "pg_catalog.timestamp":
		if notNull {
			return "time.Time"
		} else {
			return "pq.NullTime"
		}
	case "pg_catalog.timestamptz":
		if notNull {
			return "time.Time"
		} else {
			return "pq.NullTime"
		}
	case "pg_catalog.varchar":
		if notNull {
			return "string"
		} else {
			return "sql.NullString"
		}
	case "uuid":
		return "uuid.UUID"
	default:
		for name, schema := range r.Catalog.Schemas {
			if name != "public" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					return structName(enum.Name)
				}
			}
		}
		log.Printf("unknown Postgres type: %s\n", columnType)
		return "interface{}"
	}
}

func (r Result) GoQueries() []GoQuery {
	structs := r.Structs()

	var qs []GoQuery
	for _, query := range r.QueryTwos {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		// TODO: Will horribly break
		var cols []string
		for _, c := range query.Columns {
			cols = append(cols, c.Name)
		}

		gq := GoQuery{
			Cmd:          query.Cmd,
			ConstantName: lowerTitle(query.Name),
			FieldName:    lowerTitle(query.Name),
			MethodName:   query.Name,
			SQL:          strings.Replace(query.SQL, "*", strings.Join(cols, ", "), 1),
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = GoQueryValue{
				Name: p.Name,
				typ:  r.goType(p.DataType, p.NotNull),
			}
		} else if len(query.Params) > 1 {
			val := GoQueryValue{
				Emit: true,
				Name: "arg",
				Struct: &GoStruct{
					Name: gq.MethodName + "Params",
				},
			}
			for _, p := range query.Params {
				val.Struct.Fields = append(val.Struct.Fields, GoField{
					Name: structName(p.Name),
					Type: r.goType(p.DataType, p.NotNull),
					Tags: map[string]string{
						"json": p.Name,
					},
				})
			}
			gq.Arg = val
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = GoQueryValue{
				Name: c.Name,
				typ:  r.goType(c.DataType, c.NotNull),
			}
		} else if len(query.Columns) > 1 {
			var gs *GoStruct
			var emit bool

			for _, s := range structs {
				if len(s.Fields) != len(query.Columns) {
					continue
				}
				same := true
				for i, f := range s.Fields {
					c := query.Columns[i]
					sameName := f.Name == structName(c.Name)
					sameType := f.Type == r.goType(c.DataType, c.NotNull)
					if !sameName || !sameType {
						same = false
					}
				}
				if same {
					gs = &s
					break
				}
			}

			if gs == nil {
				gs = &GoStruct{
					Name: gq.MethodName + "Row",
				}
				for _, c := range query.Columns {
					gs.Fields = append(gs.Fields, GoField{
						Name: structName(c.Name),
						Type: r.goType(c.DataType, c.NotNull),
						Tags: map[string]string{"json": c.Name},
					})
				}
				emit = true
			}
			gq.Ret = GoQueryValue{
				Emit:   emit,
				Name:   "i",
				Struct: gs,
			}
		}

		qs = append(qs, gq)
	}
	sort.Slice(qs, func(i, j int) bool { return qs[i].MethodName < qs[j].MethodName })
	return qs
}

func (r Result) Records() []postgres.Table {
	used := map[string]struct{}{}
	for _, q := range r.Queries {
		used[q.ReturnType] = struct{}{}
	}
	var tables []postgres.Table
	for _, t := range r.Schema.Tables {
		if _, ok := used[t.GoName]; ok {
			tables = append(tables, t)
		}
	}
	return tables
}

func getTable(s *postgres.Schema, name string) postgres.Table {
	for _, t := range s.Tables {
		if t.Name == name {
			return t
		}
	}
	return postgres.Table{}
}

func ParseQueries(c core.Catalog, settings GenerateSettings) (*Result, error) {
	s := convert(c, settings)
	files, err := ioutil.ReadDir(settings.QueryDir)
	if err != nil {
		return nil, err
	}

	var q []Query
	var q2 []*QueryTwo

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
		r := Result{Schema: s, Settings: settings}
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
				q2 = append(q2, queryTwo)
			}
		}

		if err := parseFuncs(s, &r, source, tree); err != nil {
			return nil, err
		}
		q = append(q, r.Queries...)
	}
	return &Result{Schema: s, Queries: q, Catalog: c, QueryTwos: q2, Settings: settings}, nil
}

func parseQueryMetadata(t string) (Query, error) {
	for _, line := range strings.Split(t, "\n") {
		if !strings.HasPrefix(line, "-- name:") {
			continue
		}
		part := strings.Split(line, " ")
		return Query{
			MethodName: part[2],
			Type:       strings.TrimSpace(part[3]),
			StmtName:   lowerTitle(part[2]),
			QueryName:  lowerTitle(part[2]),
		}, nil
	}
	return Query{}, fmt.Errorf("no query metadata found")
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

func parseQuery(c core.Catalog, stmt nodes.Node, source string) (*QueryTwo, error) {
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

	return &QueryTwo{
		Cmd:     cmd,
		Name:    name,
		Params:  params,
		Columns: cols,
		SQL:     rawSQL,
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
			if postgres.IsComparisonOperator(join(n.Name, "")) {
				// TODO: Generate a name for these operations
				cols = append(cols, core.Column{Name: "_", DataType: "bool", NotNull: true})
			}

		case nodes.ColumnRef:
			parts := stringSlice(n.Fields)
			var name, alias string
			switch {
			case IsStarRef(n):
				// TODO: Disambiguate columns
				for _, t := range tables {
					for _, c := range t.Columns {
						cols = append(cols, c)
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
						cols = append(cols, c)
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
			cols = append(cols, core.Column{Name: join(n.Funcname, "."), DataType: "integer"})

		}
	}
	return cols, nil
}

func parseFuncs(s *postgres.Schema, r *Result, source string, tree pg.ParsetreeList) error {
	for _, stmt := range tree.Statements {
		if err := validateParamRef(stmt); err != nil {
			return err
		}
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			continue
		}
		switch raw.Stmt.(type) {
		case nodes.SelectStmt:
		case nodes.DeleteStmt:
		case nodes.InsertStmt:
		case nodes.UpdateStmt:
		default:
			continue
		}

		if err := validateFuncCall(raw); err != nil {
			return err
		}

		rvs := rangeVars(raw.Stmt)
		t := tableName(raw.Stmt)
		c := columnNames(s, t)

		rawSQL, _ := pluckQuery(source, raw)
		outs := findOutputs(raw.Stmt)
		refs := findParameters(raw.Stmt)

		// Super gross hack
		ctes := map[string][]outputRef{}
		{
			if selStmt, ok := raw.Stmt.(nodes.SelectStmt); ok {
				if selStmt.WithClause != nil {
					for _, item := range selStmt.WithClause.Ctes.Items {
						if cte, ok := item.(nodes.CommonTableExpr); ok {
							outs := findOutputs(cte.Ctequery)
							ctes[*cte.Ctename] = outs
						}
					}
				}
			}
		}

		tab := getTable(s, t)
		args, err := resolveRefs(s, rvs, refs)
		if err != nil {
			return err
		}

		meta, err := parseQueryMetadata(rawSQL)
		if err != nil {
			continue
		}
		meta.Table = tab
		meta.Args = args

		if len(outs) == 0 {
			meta.SQL = rawSQL
		} else if len(outs) == 1 && isStar(outs[0]) {
			meta.ReturnType = tab.GoName
			meta.ScanRecord = true
			meta.Fields = fieldsFromTable(tab)
			meta.SQL = strings.Replace(rawSQL, "*", strings.Join(c, ", "), 1)
		} else if len(outs) > 1 {
			meta.ReturnType = meta.MethodName + "Row"
			meta.ScanRecord = true
			meta.RowStruct = true
			fields, err := fieldsFromRefs(tab, ctes, outs)
			if err != nil {
				return err
			}
			meta.Fields = fields
			meta.SQL = rawSQL
		} else {
			rt, err := returnType(tab, ctes, outs)
			if err != nil {
				return err
			}
			meta.ReturnType = rt
			meta.SQL = rawSQL
		}

		r.Queries = append(r.Queries, meta)
	}
	return nil
}

func fieldsFromRefs(t postgres.Table, ctes map[string][]outputRef, refs []outputRef) ([]Field, error) {
	var f []Field
	for _, cf := range refs {
		if cf.typ != "" {
			f = append(f, Field{
				Name: strings.Title(cf.name),
				Type: cf.typ,
			})
		}
		if cf.ref != nil {
			parts := stringSlice(cf.ref.Fields)
			var name, pref string
			switch len(parts) {
			case 1:
				name = parts[0]
			case 2:
				pref = parts[0]
				name = parts[1]
			default:
				panic(fmt.Sprintf("unknown number of fields: %d", len(parts)))
			}
			var found bool
			if pref != "" && pref != t.Name {
				for _, oref := range ctes[pref] {
					if oref.name == name {
						found = true
						f = append(f, Field{
							Name: structName(pref + "_" + name),
							Type: oref.typ,
						})
					}
				}
			} else {
				for _, c := range t.Columns {
					if c.Name == name {
						found = true
						f = append(f, Field{
							Name: c.GoName,
							Type: c.GoType,
						})
					}
				}
			}
			if !found {
				return nil, Error{
					Code:    "42703",
					Message: fmt.Sprintf("column \"%s\" does not exist", name),
				}
			}
		}
	}
	return f, nil
}

func fieldsFromTable(t postgres.Table) []Field {
	var f []Field
	for _, c := range t.Columns {
		f = append(f, Field{
			Name: c.GoName,
			Type: c.GoType,
		})
	}
	return f
}

func returnType(t postgres.Table, ctes map[string][]outputRef, refs []outputRef) (string, error) {
	if len(refs) != 1 {
		return "", fmt.Errorf("too many return columns")
	}
	fields, err := fieldsFromRefs(t, ctes, refs)
	if err != nil {
		return "", err
	}
	if len(fields) == 0 {
		return "", fmt.Errorf("no fields found")
	}
	if len(fields) != 1 {
		return "", fmt.Errorf("too many fields found")
	}
	return fields[0].Type, nil
}

type outputRef struct {
	ref  *nodes.ColumnRef
	name string
	typ  string
}

type outputSearch struct {
	refs []outputRef
}

func (o *outputSearch) outs(list nodes.List) []outputRef {
	var refs []outputRef
	for _, node := range list.Items {
		res, ok := node.(nodes.ResTarget)
		if !ok {
			continue
		}
		switch n := res.Val.(type) {
		case nodes.A_Expr:
			if postgres.IsComparisonOperator(join(n.Name, "")) {
				// TODO: Generate a name for these operations
				refs = append(refs, outputRef{name: "_", typ: "bool"})
			}
		case nodes.ColumnRef:
			refs = append(refs, outputRef{ref: &n})
		case nodes.FuncCall:
			refs = append(refs, outputRef{name: join(n.Funcname, "."), typ: "int"})
		}
	}
	return refs
}

func (o *outputSearch) Visit(node nodes.Node) Visitor {
	switch n := node.(type) {
	case nodes.InsertStmt:
		o.refs = o.outs(n.ReturningList)
		return nil
	case nodes.SelectStmt:
		o.refs = o.outs(n.TargetList)
		return nil
	case nodes.UpdateStmt:
		o.refs = o.outs(n.ReturningList)
		return nil
	}
	return o
}

func findOutputs(root nodes.Node) []outputRef {
	// spew.Dump(root)
	v := &outputSearch{}
	Walk(v, root)
	return v.refs
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
	case nodes.ParamRef:
		if _, found := p.refs[n.Number]; !found {
			p.refs[n.Number] = paramRef{parent: p.parent, ref: n, rv: p.rangeVar}
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

func resolveCatalogRefs(c core.Catalog, rvs []nodes.RangeVar, args []paramRef) ([]Parameter, error) {
	typeMap := map[string]map[string]string{}
	for _, t := range c.Schemas["public"].Tables {
		typeMap[t.Name] = map[string]string{}
		for _, c := range t.Columns {
			typeMap[t.Name][c.Name] = c.DataType
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

				if typ, ok := typeMap[table][key]; ok {
					a = append(a, Parameter{Number: ref.ref.Number, Name: argName(key), DataType: typ})
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
			if typ, ok := typeMap[defaultTable][key]; ok {
				a = append(a, Parameter{Number: ref.ref.Number, Name: argName(key), DataType: typ})
			} else {
				return nil, Error{
					Code:    "42703",
					Message: fmt.Sprintf("column \"%s\" does not exist", key),
				}
			}
		case nodes.ParamRef:
			a = append(a, Parameter{Number: ref.ref.Number, Name: "_", DataType: "interface{}"})
		default:
			// return nil, fmt.Errorf("unsupported type: %T", n)
		}
	}
	return a, nil
}

func resolveRefs(s *postgres.Schema, rvs []nodes.RangeVar, args []paramRef) ([]Arg, error) {
	typeMap := map[string]map[string]string{}
	for _, t := range s.Tables {
		typeMap[t.Name] = map[string]string{}
		for _, c := range t.Columns {
			typeMap[t.Name][c.Name] = c.GoType
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

	a := []Arg{}
	for _, ref := range args {
		switch n := ref.parent.(type) {
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

				if typ, ok := typeMap[table][key]; ok {
					a = append(a, Arg{Name: argName(key), Type: typ})
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
			if typ, ok := typeMap[defaultTable][key]; ok {
				a = append(a, Arg{Name: argName(key), Type: typ})
			} else {
				return nil, Error{
					Code:    "42703",
					Message: fmt.Sprintf("column \"%s\" does not exist", key),
				}
			}
		case nodes.ParamRef:
			a = append(a, Arg{Name: "_", Type: "interface{}"})
		default:
			// return nil, fmt.Errorf("unsupported type: %T", n)
		}
	}
	return a, nil
}

func columnNames(s *postgres.Schema, table string) []string {
	cols := []string{}
	for _, t := range s.Tables {
		if t.Name != table {
			continue
		}
		for _, c := range t.Columns {
			cols = append(cols, c.Name)
		}
	}
	return cols
}

func columnType(s *postgres.Schema, settings GenerateSettings, cType string, notNull bool) string {
	for _, oride := range settings.Overrides {
		if oride.PostgresType == cType && oride.Null != notNull {
			return oride.GoType
		}
	}

	switch cType {
	case "text":
		if notNull {
			return "string"
		} else {
			return "sql.NullString"
		}
	case "serial":
		return "int"
	case "integer":
		return "int"
	case "bool":
		if notNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}
	case "jsonb":
		return "json.RawMessage"
	case "pg_catalog.bool":
		if notNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}
	case "pg_catalog.int2":
		return "uint8"
	case "pg_catalog.int4":
		return "int"
	case "pg_catalog.int8":
		return "int"
	case "pg_catalog.timestamp":
		if notNull {
			return "time.Time"
		} else {
			return "pq.NullTime"
		}
	case "pg_catalog.timestamptz":
		if notNull {
			return "time.Time"
		} else {
			return "pq.NullTime"
		}
	case "pg_catalog.varchar":
		if notNull {
			return "string"
		} else {
			return "sql.NullString"
		}
	case "uuid":
		return "uuid.UUID"
	default:
		for _, e := range s.Enums {
			if cType == e.Name {
				return e.GoName
			}
		}
		log.Printf("unknown Postgres type: %s\n", cType)
		return "interface{}"
	}
}

func tableName(n nodes.Node) string {
	switch n := n.(type) {
	case nodes.DeleteStmt:
		return *n.Relation.Relname
	case nodes.InsertStmt:
		return *n.Relation.Relname
	case nodes.SelectStmt:
		for _, item := range n.FromClause.Items {
			switch i := item.(type) {
			case nodes.RangeVar:
				return *i.Relname
			}
		}
	case nodes.UpdateStmt:
		return *n.Relation.Relname
	}
	return ""
}

var hh = `package {{.Package}}
import (
	{{- range .StdImports}}
	"{{.}}"
	{{- end}}

	{{range .PkgImports}}
	"{{.}}"
	{{- end}}
)

{{range .Enums}}
type {{.Name}} string

const (
	{{- range .Constants}}
	{{.Name}} {{.Type}} = "{{.Value}}"
	{{- end}}
)
{{end}}

{{range .Structs}}
type {{.Name}} struct { {{- range .Fields}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}

type dbtx interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db dbtx) *Queries {
	return &Queries{db: db}
}

{{if .EmitPreparedQueries}}
func Prepare(ctx context.Context, db dbtx) (*Queries, error) {
	q := Queries{db: db}
	var err error{{range .GoQueries}}
	if q.{{.FieldName}}, err = db.PrepareContext(ctx, {{.ConstantName}}); err != nil {
		return nil, err
	}
	{{- end}}
	return &q, nil
}
{{end}}

type Queries struct {
	db dbtx

    {{- if .EmitPreparedQueries}}
	tx         *sql.Tx
	{{- range .GoQueries}}
	{{.FieldName}}  *sql.Stmt
	{{- end}}
	{{- end}}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
     	{{- if .EmitPreparedQueries}}
		tx: tx,
		{{- range .GoQueries}}
		{{.FieldName}}: q.{{.FieldName}},
		{{- end}}
		{{- end}}
	}
}

{{range .GoQueries}}
const {{.ConstantName}} = {{$.Q}}{{.SQL}}
{{$.Q}}

{{if .Arg.EmitStruct}}
type {{.Arg.Type}} struct { {{- range .Arg.Struct.Fields}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}

{{if .Ret.EmitStruct}}
type {{.Ret.Type}} struct { {{- range .Ret.Struct.Fields}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}

{{if eq .Cmd ":one"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ({{.Ret.Type}}, error) {
  	{{- if $.EmitPreparedQueries}}
	var row *sql.Row
	switch {
	case q.{{.FieldName}} != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.{{.FieldName}}).QueryRowContext(ctx, {{.Arg.Params}})
	case q.{{.FieldName}} != nil:
		row = q.{{.FieldName}}.QueryRowContext(ctx, {{.Arg.Params}})
	default:
		row = q.db.QueryRowContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
	}
	{{- else}}
	row := q.db.QueryRowContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
	{{- end}}
	var {{.Ret.Name}} {{.Ret.Type}}
	err := row.Scan({{.Ret.Scan}})
	return {{.Ret.Name}}, err
}
{{end}}

{{if eq .Cmd ":many"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ([]{{.Ret.Type}}, error) {
  	{{- if $.EmitPreparedQueries}}
	var rows *sql.Rows
	var err error
	switch {
	case q.{{.FieldName}} != nil && q.tx != nil:
		rows, err = q.tx.StmtContext(ctx, q.{{.FieldName}}).QueryContext(ctx, {{.Arg.Params}})
	case q.{{.FieldName}} != nil:
		rows, err = q.{{.FieldName}}.QueryContext(ctx, {{.Arg.Params}})
	default:
		rows, err = q.db.QueryContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
	}
  	{{- else}}
	rows, err := q.db.QueryContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []{{.Ret.Type}}
	for rows.Next() {
		var {{.Ret.Name}} {{.Ret.Type}}
		if err := rows.Scan({{.Ret.Scan}}); err != nil {
			return nil, err
		}
		items = append(items, {{.Ret.Name}})
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
{{end}}

{{if eq .Cmd ":exec"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) error {
  	{{- if $.EmitPreparedQueries}}
	var err error
	switch {
	case q.{{.FieldName}} != nil && q.tx != nil:
		_, err = q.tx.StmtContext(ctx, q.{{.FieldName}}).ExecContext(ctx, {{.Arg.Params}})
	case q.{{.FieldName}} != nil:
		_, err = q.{{.FieldName}}.ExecContext(ctx, {{.Arg.Params}})
	default:
		_, err = q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
	}
  	{{- else}}
	_, err := q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	return err
}
{{end}}

{{if eq .Cmd ":execrows"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (int64, error) {
  	{{- if $.EmitPreparedQueries}}
	var result sql.Result
	var err error
	switch {
	case q.{{.FieldName}} != nil && q.tx != nil:
		result, err = q.tx.StmtContext(ctx, q.{{.FieldName}}).ExecContext(ctx, {{.Arg.Params}})
	case q.{{.FieldName}} != nil:
		result, err = q.{{.FieldName}}.ExecContext(ctx, {{.Arg.Params}})
	default:
		result, err = q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
	}
  	{{- else}}
	result, err := q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
{{end}}
{{end}}
`

func structName(name string) string {
	// if strings.HasSuffix(name, "s") {
	// 	name = name[:len(name)-1]
	// }
	out := ""
	for _, p := range strings.Split(name, "_") {
		if p == "id" {
			out += "ID"
		} else {
			out += strings.Title(p)
		}
	}
	return out
}

type Arg struct {
	Name string
	Type string
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

type tmplCtx struct {
	Q          string
	Package    string
	PkgImports []string
	StdImports []string
	Enums      []GoEnum
	Structs    []GoStruct
	GoQueries  []GoQuery

	EmitJSONTags        bool
	EmitPreparedQueries bool

	// TODO: Deprecated
	Queries []Query
	Records []postgres.Table
}

func lowerTitle(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

type TypeOverride struct {
	Package      string `json:"package"`
	PostgresType string `json:"postgres_type"`
	GoType       string `json:"go_type"`
	Null         bool   `json:"null"`
}

type GenerateSettings struct {
	SchemaDir           string         `json:"schema"`
	QueryDir            string         `json:"queries"`
	Out                 string         `json:"out"`
	Package             string         `json:"package"`
	EmitPreparedQueries bool           `json:"emit_prepared_queries"`
	EmitTags            bool           `json:"emit_tags"`
	Overrides           []TypeOverride `json:"overrides"`
}

func generate(r *Result, settings GenerateSettings) (string, error) {
	sort.Slice(r.Queries, func(i, j int) bool { return r.Queries[i].MethodName < r.Queries[j].MethodName })

	funcMap := template.FuncMap{
		"lowerTitle": lowerTitle,
	}

	pkg := "db"
	if settings.Package != "" {
		pkg = settings.Package
	}

	fileTmpl := template.Must(template.New("table").Funcs(funcMap).Parse(hh))
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	err := fileTmpl.Execute(w, tmplCtx{
		EmitPreparedQueries: settings.EmitPreparedQueries,
		EmitJSONTags:        settings.EmitTags,
		Q:                   "`",
		Queries:             r.Queries,
		GoQueries:           r.GoQueries(),
		Package:             pkg,
		Enums:               r.Enums(),
		Structs:             r.Structs(),
		Records:             r.Records(),
		StdImports:          r.StdImports(),
		PkgImports:          r.PkgImports(settings),
	})
	w.Flush()
	if err != nil {
		return "", err
	}
	code, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(b.String())
		panic(fmt.Errorf("source error: %s", err))
	}
	return string(code), nil
}
