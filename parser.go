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
		case nodes.CreateStmt:
			table := postgres.Table{
				Name:   *n.Relation.Relname,
				GoName: structName(*n.Relation.Relname),
			}
			for _, elt := range n.TableElts.Items {
				switch n := elt.(type) {
				case nodes.ColumnDef:
					// spew.Dump(n)
					// log.Printf("not null: %t", n.IsNotNull)
					table.Columns = append(table.Columns, postgres.Column{
						Name:    *n.Colname,
						GoName:  structName(*n.Colname),
						NotNull: isNotNull(n),
					})
				}
			}
			s.Tables = append(s.Tables, table)
		default:
			// spew.Dump(n)
		}
	}
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

type Query struct {
	Type       string
	MethodName string
	StmtName   string
	QueryName  string
	SQL        string
	Args       []Arg
	Table      postgres.Table
}

type Result struct {
	Schema  *postgres.Schema
	Queries []Query
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
		return &r, nil
	}
	return nil, nil
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
			t := tableName(n)
			c := columnNames(s, t)

			rawSQL, _ := pluckQuery(source, raw)
			refs := extractArgs(n)

			tab := getTable(s, t)
			r.Queries[i].Table = tab
			r.Queries[i].Args = parseArgs(tab, refs)
			r.Queries[i].SQL = strings.Replace(rawSQL, "*", strings.Join(c, ", "), 1)
		default:
			log.Printf("%T\n", n)
		}
	}
}

func extractArgs(n nodes.Node) []paramRef {
	refs := findRefs([]paramRef{}, n, nil)
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
	case nodes.RawStmt:
		r = findRefs(r, n.Stmt, nil)
	case nodes.SelectStmt:
		r = findRefs(r, n.WhereClause, nil)
		r = findRefs(r, n.LimitCount, nil)
		r = findRefs(r, n.LimitOffset, nil)
	case nodes.BoolExpr:
		for _, item := range n.Args.Items {
			r = findRefs(r, item, nil)
		}
	case nodes.A_Expr:
		r = findRefs(r, n, n.Lexpr)
		r = findRefs(r, n, n.Rexpr)
	case nodes.ParamRef:
		r = append(r, paramRef{
			parent: parent,
			ref:    n,
		})
	case nodes.ColumnRef:
	case nil:
	default:
		log.Printf("%T\n", n)
	}
	return r
}

func parseArgs(t postgres.Table, args []paramRef) []Arg {
	typeMap := map[string]string{}
	for _, c := range t.Columns {
		typeMap[c.Name] = "string"
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
					a = append(a, Arg{Name: key, Type: typ})
				} else {
					panic("unknown column: " + key)
				}
			}
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

func tableName(n nodes.SelectStmt) string {
	for _, item := range n.FromClause.Items {
		switch i := item.(type) {
		case nodes.RangeVar:
			return *i.Relname
		}
	}
	return ""
}

var hh = `package {{.Package}}
import (
	"context"
	"database/sql"
)

{{range .Schema.Tables}}
type {{.GoName}} struct { {{- range .Columns}}
  {{.GoName}} {{if .NotNull }}string{{else}}sql.NullString{{end}}
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

{{if eq .Type ":one"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) ({{.Table.GoName}}, error) {
	var row *sql.Row
	switch {
	case q.{{.StmtName}} != nil && q.tx != nil:
		row = q.tx.StmtContext(ctx, q.{{.StmtName}}).QueryRowContext(ctx, {{range .Args}}{{.Name}},{{end}})
	case q.{{.StmtName}} != nil:
		row = q.{{.StmtName}}.QueryRowContext(ctx, {{range .Args}}{{.Name}},{{end}})
	default:
		row = q.db.QueryRowContext(ctx, {{.QueryName}}, {{range .Args}}{{.Name}},{{end}})
	}
	i := {{.Table.GoName}}{}
	err := row.Scan({{range .Table.Columns}}&i.{{.GoName}},{{end}})
	return i, err
}
{{end}}

{{if eq .Type ":many"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{range .Args}}{{.Name}} {{.Type}},{{end}}) ([]{{.Table.GoName}}, error) {
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
	items := []{{.Table.GoName}}{}
	for rows.Next() {
		i := {{.Table.GoName}}{}
		if err := rows.Scan({{range .Table.Columns}}&i.{{.GoName}},{{end}}); err != nil {
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

type tmplCtx struct {
	Q       string
	Package string
	Queries []Query
	Schema  *postgres.Schema
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
		Q:       "`",
		Queries: r.Queries,
		Package: pkg,
		Schema:  r.Schema,
	})
	w.Flush()
	code, err := format.Source(b.Bytes())
	if err != nil {
		panic(fmt.Errorf("source error: %s", err))
	}
	return string(code)
}
