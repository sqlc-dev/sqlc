package strongdb

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/davecgh/go-spew/spew"
	"github.com/kyleconroy/strongdb/postgres"
	pg "github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func ParseSchmea(dir string) (*postgres.Schema, error) {
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
		tree, err := pg.Parse(string(blob))
		if err != nil {
			return nil, err
		}
		parse(&s, tree)
	}
	return &s, nil
}

func parse(s *postgres.Schema, tree pg.ParsetreeList) {
	for _, stmt := range tree.Statements {
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
				panic("could not find table " + *n.Relation.Relname)
			}

			for _, cmd := range n.Cmds.Items {
				switch cmd := cmd.(type) {
				case nodes.AlterTableCmd:
					switch cmd.Subtype {
					case nodes.AT_AddColumn:
						switch n := cmd.Def.(type) {
						case nodes.ColumnDef:
							s.Tables[idx].Columns = append(s.Tables[idx].Columns, postgres.Column{
								Name:    *n.Colname,
								Type:    join(n.TypeName.Names, "."),
								GoName:  structName(*n.Colname),
								NotNull: isNotNull(n),
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
		case nodes.CreateStmt:
			table := postgres.Table{
				Name:   *n.Relation.Relname,
				GoName: structName(*n.Relation.Relname),
			}
			for _, elt := range n.TableElts.Items {
				switch n := elt.(type) {
				case nodes.ColumnDef:
					// log.Printf("not null: %t", n.IsNotNull)
					table.Columns = append(table.Columns, postgres.Column{
						Name:    *n.Colname,
						Type:    join(n.TypeName.Names, "."),
						GoName:  structName(*n.Colname),
						NotNull: isNotNull(n),
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
					panic("could not find table " + *n.Relation.Relname)
				}
				s.Tables[idx].Name = *n.Newname
				s.Tables[idx].GoName = structName(*n.Newname)
			}
		default:
			// spew.Dump(n)
		}
	}
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

func isStar(n nodes.ColumnRef) bool {
	if len(n.Fields.Items) != 1 {
		return false
	}
	_, aStar := n.Fields.Items[0].(nodes.A_Star)
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

func (r Result) UsesTime() bool {
	for _, table := range r.Records() {
		for _, c := range table.Columns {
			if c.GoType() == "time.Time" {
				return true
			}
		}
	}
	return false
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
		blob, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		r := Result{Schema: s, Queries: parseQueries(blob)}
		tree, err := pg.Parse(string(blob))
		if err != nil {
			return nil, err
		}
		parseFuncs(s, &r, string(blob), tree)
		q = append(q, r.Queries...)
	}
	return &Result{Schema: s, Queries: q}, nil
}

func parseQueries(t []byte) []Query {
	q := []Query{}
	for _, line := range strings.Split(string(t), "\n") {
		if !strings.HasPrefix(line, "-- name:") {
			continue
		}
		part := strings.Split(line, " ")
		q = append(q, Query{
			MethodName: part[2],
			Type:       strings.TrimSpace(part[3]),
			StmtName:   lowerTitle(part[2]),
			QueryName:  lowerTitle(part[2]),
		})
	}
	return q
}

func pluckQuery(source string, n nodes.RawStmt) (string, error) {
	// TODO: Bounds checking
	head := n.StmtLocation
	tail := n.StmtLocation + n.StmtLen
	return strings.TrimSpace(source[head:tail]), nil
}

func parseFuncs(s *postgres.Schema, r *Result, source string, tree pg.ParsetreeList) {
	for i, stmt := range tree.Statements {
		raw, ok := stmt.(nodes.RawStmt)
		if !ok {
			continue
		}
		switch n := raw.Stmt.(type) {
		case nodes.SelectStmt:
		case nodes.DeleteStmt:
		case nodes.InsertStmt:
		case nodes.UpdateStmt:
		default:
			log.Printf("%T\n", n)
			continue
		}

		t := tableName(raw.Stmt)
		c := columnNames(s, t)

		rawSQL, _ := pluckQuery(source, raw)
		refs := extractArgs(raw.Stmt)
		outs := findOutputs(nil, raw.Stmt)

		tab := getTable(s, t)
		r.Queries[i].Table = tab
		r.Queries[i].Args = parseArgs(tab, refs)

		if len(outs) == 0 {
			r.Queries[i].SQL = rawSQL
		} else if len(outs) == 1 && isStar(outs[0]) {
			r.Queries[i].ReturnType = tab.GoName
			r.Queries[i].ScanRecord = true
			r.Queries[i].Fields = fieldsFromTable(tab)
			r.Queries[i].SQL = strings.Replace(rawSQL, "*", strings.Join(c, ", "), 1)
		} else if len(outs) > 1 {
			r.Queries[i].ReturnType = r.Queries[i].MethodName + "Row"
			r.Queries[i].ScanRecord = true
			r.Queries[i].RowStruct = true
			r.Queries[i].Fields = fieldsFromRefs(tab, outs)
			r.Queries[i].SQL = rawSQL
		} else {
			r.Queries[i].ReturnType = returnType(tab, outs)
			r.Queries[i].SQL = rawSQL
		}
	}
}

func fieldsFromRefs(t postgres.Table, refs []nodes.ColumnRef) []Field {
	var f []Field
	for _, cf := range refs {
		name := join(cf.Fields, ".")
		for _, c := range t.Columns {
			if c.Name == name {
				f = append(f, Field{
					Name: c.GoName,
					Type: c.GoType(),
				})
			}
		}
	}
	return f
}

func fieldsFromTable(t postgres.Table) []Field {
	var f []Field
	for _, c := range t.Columns {
		f = append(f, Field{
			Name: c.GoName,
			Type: c.GoType(),
		})
	}
	return f
}

func returnType(t postgres.Table, refs []nodes.ColumnRef) string {
	if len(refs) != 1 {
		// panic("too many return columns")
		return "interface{}"
	}
	name := join(refs[0].Fields, ".")
	for _, c := range t.Columns {
		if c.Name == name {
			return c.GoType()
		}
	}
	return "interface{}"
}

func extractArgs(n nodes.Node) []paramRef {
	allrefs := findRefs([]paramRef{}, n, nil)
	refs := make([]paramRef, 0)
	seen := map[int]struct{}{}
	for _, r := range allrefs {
		if _, ok := seen[r.ref.Number]; ok {
			continue
		}
		refs = append(refs, r)
		seen[r.ref.Number] = struct{}{}
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].ref.Number < refs[j].ref.Number })
	return refs
}

type paramRef struct {
	parent nodes.Node
	ref    nodes.ParamRef
}

func findRefs(r []paramRef, parent, n nodes.Node) []paramRef {
	if n == nil {
		n = parent
	}
	switch n := n.(type) {
	case nodes.A_Expr:
		r = findRefs(r, n, n.Lexpr)
		r = findRefs(r, n, n.Rexpr)
	case nodes.ColumnRef:
	case nodes.BoolExpr:
		r = findRefs(r, n.Args, nil)
	case nodes.DeleteStmt:
		r = findRefs(r, n.WhereClause, nil)
	case nodes.FuncCall:
	case nodes.InsertStmt:
		switch s := n.SelectStmt.(type) {
		case nodes.SelectStmt:
			for _, vl := range s.ValuesLists {
				for i, v := range vl {
					// TODO: Index error
					r = findRefs(r, n.Cols.Items[i], v)
				}
			}
		}
	case nodes.List:
		for _, item := range n.Items {
			r = findRefs(r, item, nil)
		}
	case nodes.ParamRef:
		r = append(r, paramRef{
			parent: parent,
			ref:    n,
		})
	case nodes.RawStmt:
		r = findRefs(r, n.Stmt, nil)
	case nodes.ResTarget:
		r = findRefs(r, n, n.Val)
	case nodes.SelectStmt:
		r = findRefs(r, n.WhereClause, nil)
		r = findRefs(r, n.LimitCount, nil)
		r = findRefs(r, n.LimitOffset, nil)
	case nodes.UpdateStmt:
		r = findRefs(r, n.TargetList, nil)
		r = findRefs(r, n.WhereClause, nil)
	case nil:
	default:
		log.Printf("%T\n", n)
	}
	return r
}

func findOutputs(r []nodes.ColumnRef, n nodes.Node) []nodes.ColumnRef {
	switch n := n.(type) {
	case nodes.ColumnRef:
		r = append(r, n)
	case nodes.DeleteStmt:
		r = findOutputs(r, n.ReturningList)
	case nodes.FuncCall:
		// join(n.Funcname.List, ".")
		spew.Dump(n)
	case nodes.InsertStmt:
		r = findOutputs(r, n.ReturningList)
	case nodes.List:
		for _, i := range n.Items {
			r = findOutputs(r, i)
		}
	case nodes.RawStmt:
		r = findOutputs(r, n.Stmt)
	case nodes.ResTarget:
		r = findOutputs(r, n.Val)
	case nodes.SelectStmt:
		r = findOutputs(r, n.TargetList)
	case nodes.UpdateStmt:
		r = findOutputs(r, n.ReturningList)
	case nil:
	default:
		log.Printf("%T\n", n)
	}
	return r
}

func parseArgs(t postgres.Table, args []paramRef) []Arg {
	typeMap := map[string]string{}
	for _, c := range t.Columns {
		typeMap[c.Name] = c.GoType()
	}
	a := []Arg{}
	for _, ref := range args {
		switch n := ref.parent.(type) {
		case nodes.A_Expr:
			switch n := n.Lexpr.(type) {
			case nodes.ColumnRef:
				key := ""
				for _, n := range n.Fields.Items {
					switch n := n.(type) {
					case nodes.String:
						key += n.Str
					}
				}
				if typ, ok := typeMap[key]; ok {
					a = append(a, Arg{Name: argName(key), Type: typ})
				} else {
					panic("unknown column: " + key)
				}
			}
		case nodes.ResTarget:
			key := *n.Name
			if typ, ok := typeMap[key]; ok {
				a = append(a, Arg{Name: argName(key), Type: typ})
			} else {
				panic("unknown column: " + key)
			}
		case nodes.ParamRef:
			a = append(a, Arg{Name: "_", Type: "interface{}"})
		default:
			panic(fmt.Sprintf("unsupported type: %T", n))
		}
	}
	return a
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
	"context"
	"database/sql"
	{{if .ImportTime}}"time"{{end}}
)

{{range .Records}}
type {{.GoName}} struct { {{- range .Columns}}
  {{.GoName}} {{.GoType}}
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

func Prepare(ctx context.Context, db dbtx) (*Queries, error) {
	q := Queries{db: db}
	var err error{{range .Queries}}
	if q.{{.StmtName}}, err = db.PrepareContext(ctx, {{.QueryName}}); err != nil {
		return nil, err
	}
	{{- end}}
	return &q, nil
}

type Queries struct {
	db dbtx

	tx         *sql.Tx
	{{- range .Queries}}
	{{.StmtName}}  *sql.Stmt
	{{- end}}
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		tx: tx,
		db: tx,
		{{- range .Queries}}
		{{.StmtName}}: q.{{.StmtName}},
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
	var row *sql.Row
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.{{.StmtName}}).QueryRowContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		row = q.{{.StmtName}}.QueryRowContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		row = q.db.QueryRowContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
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
	var err error
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		_, err = q.tx.StmtContext(ctx, q.{{.StmtName}}).ExecContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		_, err = q.{{.StmtName}}.ExecContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		_, err = q.db.ExecContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
	return err
}
{{end}}

{{if eq .Type ":execrows"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) (int64, error) {
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
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
{{end}}

{{end}}
`

func structName(name string) string {
	if strings.HasSuffix(name, "s") {
		name = name[:len(name)-1]
	}
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
	Queries    []Query
	Schema     *postgres.Schema
	Records    []postgres.Table
	ImportTime bool
}

func lowerTitle(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func generate(r *Result, pkg string) string {
	sort.Slice(r.Queries, func(i, j int) bool { return r.Queries[i].MethodName < r.Queries[j].MethodName })

	funcMap := template.FuncMap{
		"lowerTitle": lowerTitle,
	}

	fileTmpl := template.Must(template.New("table").Funcs(funcMap).Parse(hh))
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	fileTmpl.Execute(w, tmplCtx{
		Q:          "`",
		Queries:    r.Queries,
		Package:    pkg,
		Schema:     r.Schema,
		Records:    r.Records(),
		ImportTime: r.UsesTime(),
	})
	w.Flush()
	code, err := format.Source(b.Bytes())
	if err != nil {
		fmt.Println(b.String())
		panic(fmt.Errorf("source error: %s", err))
	}
	return string(code)
}
