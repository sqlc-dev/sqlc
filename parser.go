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

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/inflection"
	"github.com/kyleconroy/dinosql/postgres"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func parseSQL(in string) (*Result, error) {
	s := postgres.Schema{}
	tree, err := pg.Parse(in)
	if err != nil {
		return nil, err
	}
	if err := parse(&s, tree, GenerateSettings{}); err != nil {
		return nil, err
	}

	var q []Query
	r := Result{Schema: &s}
	if err := parseFuncs(&s, &r, in, tree); err != nil {
		return nil, err
	}
	q = append(q, r.Queries...)

	return &Result{Schema: &s, Queries: q}, nil
}

func ParseSchmea(dir string, settings GenerateSettings) (*postgres.Schema, error) {
	// Keep the import around
	if false {
		spew.Dump(dir)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	s := postgres.Schema{}
	for _, f := range files {
		blob, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		contents := RemoveGooseRollback(string(blob))
		tree, err := pg.Parse(contents)
		if err != nil {
			return nil, err
		}
		if err := parse(&s, tree, settings); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

func parse(s *postgres.Schema, tree pg.ParsetreeList, settings GenerateSettings) error {
	for _, stmt := range tree.Statements {
		if err := validateFuncCall(stmt); err != nil {
			return err
		}
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			continue
		}
		switch n := raw.Stmt.(type) {
		case nodes.AlterTableStmt:
			idx := -1
			for i, table := range s.Tables {
				if table.Name == *n.Relation.Relname {
					idx = i
				}
			}
			if idx < 0 {
				return Error{
					Code:    "42P01",
					Message: fmt.Sprintf("relation \"%s\" does not exist", *n.Relation.Relname),
				}
			}

			for _, cmd := range n.Cmds.Items {
				switch cmd := cmd.(type) {
				case nodes.AlterTableCmd:
					switch cmd.Subtype {
					case nodes.AT_AddColumn:
						switch n := cmd.Def.(type) {
						case nodes.ColumnDef:
							ctype := join(n.TypeName.Names, ".")
							notNull := isNotNull(n)
							s.Tables[idx].Columns = append(s.Tables[idx].Columns, postgres.Column{
								Name:    *n.Colname,
								Type:    ctype,
								NotNull: notNull,
								GoName:  structName(*n.Colname),
								GoType:  columnType(s, settings, ctype, notNull),
							})
						}
					case nodes.AT_DropColumn:
						for i, c := range s.Tables[idx].Columns {
							if c.Name == *cmd.Name {
								s.Tables[idx].Columns = append(s.Tables[idx].Columns[:i], s.Tables[idx].Columns[i+1:]...)
							}
						}
					}
				}
			}
		case nodes.CreateEnumStmt:
			vals := []string{}
			for _, item := range n.Vals.Items {
				if n, ok := item.(nodes.String); ok {
					vals = append(vals, n.Str)
				}
			}
			s.Enums = append(s.Enums, postgres.Enum{
				Name:   join(n.TypeName, "."),
				GoName: structName(join(n.TypeName, ".")),
				Vals:   vals,
			})
		case nodes.CreateStmt:
			table := postgres.Table{
				Name:   *n.Relation.Relname,
				GoName: inflection.Singular(structName(*n.Relation.Relname)),
			}
			for _, elt := range n.TableElts.Items {
				switch n := elt.(type) {
				case nodes.ColumnDef:
					// log.Printf("not null: %t", n.IsNotNull)
					ctype := join(n.TypeName.Names, ".")
					notNull := isNotNull(n)
					table.Columns = append(table.Columns, postgres.Column{
						Name:    *n.Colname,
						Type:    ctype,
						NotNull: notNull,
						GoName:  structName(*n.Colname),
						GoType:  columnType(s, settings, ctype, notNull),
					})
				}
			}
			s.Tables = append(s.Tables, table)
		case nodes.RenameStmt:
			switch n.RenameType {
			case nodes.OBJECT_TABLE:
				idx := -1
				for i, table := range s.Tables {
					if table.Name == *n.Relation.Relname {
						idx = i
					}
				}
				if idx < 0 {
					return Error{
						Code:    "42P01",
						Message: fmt.Sprintf("relation \"%s\" does not exist", *n.Relation.Relname),
					}
				}
				s.Tables[idx].Name = *n.Newname
				s.Tables[idx].GoName = inflection.Singular(structName(*n.Newname))
			}
		default:
			// spew.Dump(n)
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

func isNotNull(n nodes.ColumnDef) bool {
	if n.IsNotNull {
		return true
	}
	for _, c := range n.Constraints.Items {
		switch n := c.(type) {
		case nodes.Constraint:
			if n.Contype == nodes.CONSTR_NOTNULL {
				return true
			}
			if n.Contype == nodes.CONSTR_PRIMARY {
				return true
			}
		}
	}
	return false
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

type Result struct {
	Schema  *postgres.Schema
	Queries []Query
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

func (r Result) UsesType(typ string) bool {
	for _, table := range r.Records() {
		for _, c := range table.Columns {
			if c.GoType == typ {
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

func getTable(s *postgres.Schema, name string) postgres.Table {
	for _, t := range s.Tables {
		if t.Name == name {
			return t
		}
	}
	return postgres.Table{}
}

func ParseQueries(s *postgres.Schema, dir string) (*Result, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var q []Query
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".sql") {
			continue
		}
		blob, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		r := Result{Schema: s}
		tree, err := pg.Parse(string(blob))
		if err != nil {
			return nil, err
		}
		if err := parseFuncs(s, &r, string(blob), tree); err != nil {
			return nil, err
		}
		q = append(q, r.Queries...)
	}
	return &Result{Schema: s, Queries: q}, nil
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

type appender struct {
	refs []outputRef
}

func (a *appender) Visit(node nodes.Node) Visitor {
	res, ok := node.(nodes.ResTarget)
	if !ok {
		return a
	}
	switch n := res.Val.(type) {
	case nodes.A_Expr:
		if postgres.IsComparisonOperator(join(n.Name, "")) {
			// TODO: Generate a name for these operations
			a.refs = append(a.refs, outputRef{name: "_", typ: "bool"})
		}
	case nodes.ColumnRef:
		a.refs = append(a.refs, outputRef{ref: &n})
	case nodes.FuncCall:
		a.refs = append(a.refs, outputRef{name: join(n.Funcname, "."), typ: "int"})
	}
	return nil
}

type outputSearch struct {
	a *appender
}

func (o *outputSearch) Visit(node nodes.Node) Visitor {
	switch n := node.(type) {
	case nodes.InsertStmt:
		Walk(o.a, n.ReturningList)
		return nil
	case nodes.SelectStmt:
		Walk(o.a, n.TargetList)
		return nil
	case nodes.UpdateStmt:
		Walk(o.a, n.ReturningList)
		return nil
	}
	return o
}

func findOutputs(root nodes.Node) []outputRef {
	// spew.Dump(root)
	v := &outputSearch{&appender{}}
	Walk(v, root)
	return v.a.refs
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

{{range .Schema.Enums}}
type {{.GoName}} string

const (
	{{- range .Constants}}
	{{.Name}} {{.Type}} = "{{.Value}}"
	{{- end}}
)
{{end}}

{{range .Records}}
type {{.GoName}} struct { {{- range .Columns}}
  {{.GoName}} {{.GoType}} {{if $.EmitTags}}{{$.Q}}json:"{{.Name}}"{{$.Q}}{{end}}
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

{{if .PrepareSupport}}
func Prepare(ctx context.Context, db dbtx) (*Queries, error) {
	q := Queries{db: db}
	var err error{{range .Queries}}
	if q.{{.StmtName}}, err = db.PrepareContext(ctx, {{.QueryName}}); err != nil {
		return nil, err
	}
	{{- end}}
	return &q, nil
}
{{end}}

type Queries struct {
	db dbtx

    {{- if .PrepareSupport}}
	tx         *sql.Tx
	{{- range .Queries}}
	{{.StmtName}}  *sql.Stmt
	{{- end}}
	{{- end}}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
     	{{- if .PrepareSupport}}
		tx: tx,
		{{- range .Queries}}
		{{.StmtName}}: q.{{.StmtName}},
		{{- end}}
		{{- end}}
	}
}

{{range .Queries}}
const {{.QueryName}} = {{$.Q}}{{.SQL}}
{{$.Q}}

{{if .RowStruct}}
type {{.MethodName}}Row struct { {{- range .Fields}}
  {{.Name}} {{.Type}}
  {{- end}}
}
{{end}}

{{if eq .Type ":one"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) ({{.ReturnType}}, error) {
  	{{- if $.PrepareSupport}}
	var row *sql.Row
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.{{.StmtName}}).QueryRowContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		row = q.{{.StmtName}}.QueryRowContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		row = q.db.QueryRowContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
	{{- else}}
	row := q.db.QueryRowContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	{{- end}}
	var i {{.ReturnType}}
	{{- if .ScanRecord}}
	err := row.Scan({{range .Fields}}&i.{{.Name}},{{end}})
	{{- else}}
	err := row.Scan(&i)
	{{- end}}
	return i, err
}
{{end}}

{{if eq .Type ":many"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) ([]{{.ReturnType}}, error) {
  	{{- if $.PrepareSupport}}
	var rows *sql.Rows
	var err error
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		rows, err = q.tx.StmtContext(ctx, q.{{.StmtName}}).QueryContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		rows, err = q.{{.StmtName}}.QueryContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		rows, err = q.db.QueryContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
  	{{- else}}
	rows, err := q.db.QueryContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
  	{{- end}}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []{{.ReturnType}}{}
	for rows.Next() {
		var i {{.ReturnType}}
		{{- if .ScanRecord}}
		if err := rows.Scan({{range .Fields}}&i.{{.Name}},{{end}}); err != nil {
		{{- else}}
		if err := rows.Scan(&i); err != nil {
		{{- end}}
			return nil, err
		}
		items = append(items, i)
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

{{if eq .Type ":exec"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) error {
  	{{- if $.PrepareSupport}}
	var err error
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		_, err = q.tx.StmtContext(ctx, q.{{.StmtName}}).ExecContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		_, err = q.{{.StmtName}}.ExecContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		_, err = q.db.ExecContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
  	{{- else}}
	_, err := q.db.ExecContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
  	{{- end}}
	return err
}
{{end}}

{{if eq .Type ":execrows"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) (int64, error) {
  	{{- if $.PrepareSupport}}
	var result sql.Result
	var err error
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		result, err = q.tx.StmtContext(ctx, q.{{.StmtName}}).ExecContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		result, err = q.{{.StmtName}}.ExecContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		result, err = q.db.ExecContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
	{{- else}}
	result, err := q.db.ExecContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
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
	Q              string
	Package        string
	Queries        []Query
	Schema         *postgres.Schema
	Records        []postgres.Table
	StdImports     []string
	PkgImports     []string
	PrepareSupport bool
	EmitTags       bool
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

func generate(r *Result, settings GenerateSettings) string {
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
	fileTmpl.Execute(w, tmplCtx{
		PrepareSupport: settings.EmitPreparedQueries,
		EmitTags:       settings.EmitTags,
		Q:              "`",
		Queries:        r.Queries,
		Package:        pkg,
		Schema:         r.Schema,
		Records:        r.Records(),
		StdImports:     r.StdImports(),
		PkgImports:     r.PkgImports(settings),
	})
	w.Flush()
	code, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(b.String())
		panic(fmt.Errorf("source error: %s", err))
	}
	return string(code)
}
