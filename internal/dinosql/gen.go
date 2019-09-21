package dinosql

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"unicode"

	core "github.com/kyleconroy/sqlc/internal/pg"

	"github.com/jinzhu/inflection"
)

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
		if strings.HasPrefix(v.typ, "[]") {
			out = append(out, "pq.Array("+v.Name+")")
		} else {
			out = append(out, v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") {
				out = append(out, "pq.Array("+v.Name+"."+f.Name+")")
			} else {
				out = append(out, v.Name+"."+f.Name)
			}
		}
	}
	if len(out) <= 3 {
		return strings.Join(out, ",")
	}
	out = append(out, "")
	return "\n" + strings.Join(out, ",\n")
}

func (v GoQueryValue) Scan() string {
	var out []string
	if v.Struct == nil {
		if strings.HasPrefix(v.typ, "[]") {
			out = append(out, "pq.Array(&"+v.Name+")")
		} else {
			out = append(out, "&"+v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") {
				out = append(out, "pq.Array(&"+v.Name+"."+f.Name+")")
			} else {
				out = append(out, "&"+v.Name+"."+f.Name)
			}
		}
	}
	if len(out) <= 3 {
		return strings.Join(out, ",")
	}
	out = append(out, "")
	return "\n" + strings.Join(out, ",\n")
}

// A struct used to generate methods and fields on the Queries struct
type GoQuery struct {
	Cmd          string
	MethodName   string
	FieldName    string
	ConstantName string
	SQL          string
	SourceName   string
	Ret          GoQueryValue
	Arg          GoQueryValue
}

func (r Result) UsesType(typ string) bool {
	for _, strct := range r.Structs() {
		for _, f := range strct.Fields {
			if strings.HasPrefix(f.Type, typ) {
				return true
			}
		}
	}
	return false
}

func (r Result) UsesArrays() bool {
	for _, strct := range r.Structs() {
		for _, f := range strct.Fields {
			if strings.HasPrefix(f.Type, "[]") {
				return true
			}
		}
	}
	return false
}

func (r Result) Imports(filename string) [][]string {
	if filename == "db.go" {
		return [][]string{
			[]string{"context", "database/sql"},
		}
	}

	if filename == "models.go" {
		return r.ModelImports()
	}

	return r.QueryImports(filename)
}

func (r Result) ModelImports() [][]string {

	var std []string
	if r.UsesType("sql.Null") {
		std = append(std, "database/sql")
	}
	if r.UsesType("json.RawMessage") {
		std = append(std, "encoding/json")
	}
	if r.UsesType("time.Time") {
		std = append(std, "time")
	}

	var pkg []string
	if r.UsesType("pq.NullTime") {
		pkg = append(pkg, "github.com/lib/pq")
	}
	if r.UsesType("uuid.UUID") {
		pkg = append(pkg, "github.com/google/uuid")
	}

	// Custom imports
	overrideTypes := map[string]string{}
	for _, o := range r.Settings.Overrides {
		overrideTypes[o.GoType] = o.Package
	}
	for goType, importPath := range overrideTypes {
		if r.UsesType(goType) {
			pkg = append(pkg, importPath)
		}
	}

	return [][]string{std, pkg}
}

func (r Result) QueryImports(filename string) [][]string {
	// for _, strct := range r.Structs() {
	// 	for _, f := range strct.Fields {
	// 		if strings.HasPrefix(f.Type, "[]") {
	// 			return true
	// 		}
	// 	}
	// }
	var gq []GoQuery
	for _, query := range r.GoQueries() {
		if query.SourceName == filename {
			gq = append(gq, query)
		}
	}

	uses := func(name string) bool {
		for _, q := range gq {
			if !q.Ret.isEmpty() {
				if q.Ret.EmitStruct() {
					for _, f := range q.Ret.Struct.Fields {
						if strings.HasPrefix(f.Type, name) {
							return true
						}
					}
				}
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.EmitStruct() {
					for _, f := range q.Arg.Struct.Fields {
						if strings.HasPrefix(f.Type, name) {
							return true
						}
					}
				}
				if strings.HasPrefix(q.Arg.Type(), name) {
					return true
				}
			}
		}
		return false
	}

	sliceScan := func() bool {
		for _, q := range gq {
			if !q.Ret.isEmpty() {
				if q.Ret.IsStruct() {
					for _, f := range q.Ret.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Ret.Type(), "[]") {
						return true
					}
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.IsStruct() {
					for _, f := range q.Arg.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Arg.Type(), "[]") {
						return true
					}
				}
			}
		}
		return false
	}

	std := []string{"context"}
	if uses("sql.Null") {
		std = append(std, "database/sql")
	}
	if uses("json.RawMessage") {
		std = append(std, "encoding/json")
	}
	if uses("time.Time") {
		std = append(std, "time")
	}

	var pkg []string
	if sliceScan() {
		pkg = append(pkg, "github.com/lib/pq")
	}
	if uses("pq.NullTime") {
		pkg = append(pkg, "github.com/lib/pq")
	}
	if uses("uuid.UUID") {
		pkg = append(pkg, "github.com/google/uuid")
	}

	// Custom imports
	overrideTypes := map[string]string{}
	for _, o := range r.Settings.Overrides {
		overrideTypes[o.GoType] = o.Package
	}
	for goType, importPath := range overrideTypes {
		if uses(goType) {
			pkg = append(pkg, importPath)
		}
	}

	return [][]string{std, pkg}
}

func (r Result) Enums() []GoEnum {
	var enums []GoEnum
	for name, schema := range r.Catalog.Schemas {
		if name != "public" {
			continue
		}
		for _, enum := range schema.Enums {
			e := GoEnum{
				Name: r.structName(enum.Name),
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

func (r Result) structName(name string) string {
	if rename, _ := r.Settings.Rename[name]; rename != "" {
		return rename
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

func (r Result) Structs() []GoStruct {
	var structs []GoStruct
	for name, schema := range r.Catalog.Schemas {
		if name == "pg_catalog" {
			continue
		}
		for _, table := range schema.Tables {
			s := GoStruct{
				Name: inflection.Singular(r.structName(table.Name)),
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, GoField{
					Name: r.structName(column.Name),
					Type: r.goType(column),
					Tags: map[string]string{"json:": column.Name},
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

func (r Result) goType(col core.Column) string {
	typ := r.goInnerType(col.DataType, col.NotNull || col.IsArray)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func (r Result) goInnerType(columnType string, notNull bool) string {
	for _, oride := range r.Settings.Overrides {
		if oride.PostgresType == columnType && oride.Null != notNull {
			return oride.GoType
		}
	}

	switch columnType {

	case "serial", "pg_catalog.serial4":
		return "int32"

	case "bigserial", "pg_catalog.serial8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64" // unnecessay else

	case "smallserial", "pg_catalog.serial2":
		return "int16"

	case "integer", "int", "pg_catalog.int4":
		return "int32"

	case "bigint", "pg_catalog.int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64" // unnecessary else

	case "smallint", "pg_catalog.int2":
		return "int16"

	case "float", "double precision", "pg_catalog.float8":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64" // unnecessary else

	case "real", "pg_catalog.float4":
		if notNull {
			return "float32"
		} // unnecessary else
		return "sql.NullFloat64" //IMPORTANT: Change to sql.NullFloat32 after updating the go.mod file

	case "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool" // unnecessary else

	case "jsonb":
		return "json.RawMessage"

	case "pg_catalog.timestamp", "pg_catalog.timestamptz":
		if notNull {
			return "time.Time"
		}
		return "pq.NullTime" // unnecessary else

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar":
		if notNull {
			return "string"
		}
		return "sql.NullString" // unnecessary else

	case "uuid":
		return "uuid.UUID"

	case "any":
		return "interface{}"

	default:
		for name, schema := range r.Catalog.Schemas {
			if name != "public" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					return r.structName(enum.Name)
				}
			}
		}
		log.Printf("unknown Postgres type: %s\n", columnType)
		return "interface{}"
	}
}

// It's possible that this method will generate duplicate JSON tag values
//
//   Columns: count, count,   count_2
//    Fields: Count, Count_2, Count2
// JSON tags: count, count_2, count_2
//
// This is unlikely to happen, so don't fix it yet
func (r Result) columnsToStruct(name string, columns []core.Column) *GoStruct {
	gs := GoStruct{
		Name: name,
	}
	seen := map[string]int{}
	for i, c := range columns {
		tagName := c.Name
		fieldName := r.structName(columnName(c, i))
		if v := seen[c.Name]; v > 0 {
			tagName = fmt.Sprintf("%s_%d", tagName, v+1)
			fieldName = fmt.Sprintf("%s_%d", fieldName, v+1)
		}
		gs.Fields = append(gs.Fields, GoField{
			Name: fieldName,
			Type: r.goType(c),
			Tags: map[string]string{"json:": tagName},
		})
		seen[c.Name] += 1
	}
	return &gs
}

func paramName(p Parameter) string {
	if p.Column.Name != "" {
		return p.Column.Name
	}
	return fmt.Sprintf("dollar_%d", p.Number)
}

func columnName(c core.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func (r Result) GoQueries() []GoQuery {
	structs := r.Structs()

	var qs []GoQuery
	for _, query := range r.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		code := query.SQL

		// TODO: Will horribly break sometimes
		if query.NeedsEdit {
			var cols []string
			find := "*"
			for _, c := range query.Columns {
				if c.Scope != "" {
					find = c.Scope + ".*"
				}
				name := c.Name
				if c.Scope != "" {
					name = c.Scope + "." + name
				}
				cols = append(cols, name)
			}
			code = strings.Replace(query.SQL, find, strings.Join(cols, ", "), 1)
		}

		gq := GoQuery{
			Cmd:          query.Cmd,
			ConstantName: lowerTitle(query.Name),
			FieldName:    lowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			SourceName:   query.Filename,
			SQL:          code,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = GoQueryValue{
				Name: paramName(p),
				typ:  r.goType(p.Column),
			}
		} else if len(query.Params) > 1 {
			var cols []core.Column
			for _, p := range query.Params {
				cols = append(cols, p.Column)
			}
			gq.Arg = GoQueryValue{
				Emit:   true,
				Name:   "arg",
				Struct: r.columnsToStruct(gq.MethodName+"Params", cols),
			}
		}

		if len(query.Columns) == 1 {
			c := query.Columns[0]
			gq.Ret = GoQueryValue{
				Name: columnName(c, 0),
				typ:  r.goType(c),
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
					sameName := f.Name == r.structName(columnName(c, i))
					sameType := f.Type == r.goType(c)
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
				gs = r.columnsToStruct(gq.MethodName+"Row", query.Columns)
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

var dbTmpl = `// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

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

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Row) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
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
`
var modelsTmpl = `// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
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
`

var sqlTmpl = `// Code generated by sqlc. DO NOT EDIT.
// source: {{.SourceName}}

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{range .GoQueries}}
{{if eq .SourceName $.SourceName}}
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
	row := q.queryRow(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
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
	rows, err := q.query(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
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
	_, err := q.exec(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
  	{{- else}}
	_, err := q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
	return err
}
{{end}}

{{if eq .Cmd ":execrows"}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (int64, error) {
  	{{- if $.EmitPreparedQueries}}
	result, err := q.exec(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
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
{{end}}
`

type tmplCtx struct {
	Q         string
	Package   string
	Enums     []GoEnum
	Structs   []GoStruct
	GoQueries []GoQuery
	Settings  GenerateSettings

	// TODO: Race conditions
	SourceName string

	EmitJSONTags        bool
	EmitPreparedQueries bool
}

func lowerTitle(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func Generate(r *Result, global GenerateSettings, settings PackageSettings) (map[string]string, error) {
	funcMap := template.FuncMap{
		"lowerTitle": lowerTitle,
		"imports":    r.Imports,
	}

	pkg := settings.Name
	if pkg == "" {
		pkg = filepath.Base(settings.Path)
	}

	dbFile := template.Must(template.New("table").Funcs(funcMap).Parse(dbTmpl))
	modelsFile := template.Must(template.New("table").Funcs(funcMap).Parse(modelsTmpl))
	sqlFile := template.Must(template.New("table").Funcs(funcMap).Parse(sqlTmpl))

	tctx := tmplCtx{
		Settings:            global,
		EmitPreparedQueries: settings.EmitPreparedQueries,
		EmitJSONTags:        settings.EmitJSONTags,
		Q:                   "`",
		Package:             pkg,
		GoQueries:           r.GoQueries(),
		Enums:               r.Enums(),
		Structs:             r.Structs(),
	}

	output := map[string]string{}

	execute := func(name string, t *template.Template) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := t.Execute(w, tctx)
		w.Flush()
		if err != nil {
			return err
		}
		code, err := format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return fmt.Errorf("source error: %s", err)
		}
		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}
		output[name] = string(code)
		return nil
	}

	if err := execute("db.go", dbFile); err != nil {
		return nil, err
	}
	if err := execute("models.go", modelsFile); err != nil {
		return nil, err
	}

	files := map[string]struct{}{}
	for _, gq := range r.GoQueries() {
		files[gq.SourceName] = struct{}{}
	}

	for source, _ := range files {
		if err := execute(source, sqlFile); err != nil {
			return nil, err
		}
	}
	return output, nil
}
