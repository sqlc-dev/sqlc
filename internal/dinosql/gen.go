package dinosql

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"sort"
	"strings"
	"text/template"
	"unicode"

	core "github.com/kyleconroy/dinosql/internal/pg"

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
	return strings.Join(out, ",")
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
	if r.UsesArrays() {
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
		if name != "public" {
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
		return "int64"

	case "smallserial", "pg_catalog.serial2":
		return "int16"

	case "integer", "int", "pg_catalog.int4":
		return "int32"

	case "bigint", "pg_catalog.int8":
		return "int64"

	case "smallint", "pg_catalog.int2":
		return "int16"

	case "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		} else {
			return "sql.NullBool"
		}

	case "jsonb":
		return "json.RawMessage"

	case "pg_catalog.timestamp", "pg_catalog.timestamptz":
		if notNull {
			return "time.Time"
		} else {
			return "pq.NullTime"
		}

	case "text", "pg_catalog.varchar":
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
			Tags: map[string]string{"json": tagName},
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
			for _, c := range query.Columns {
				cols = append(cols, c.Name)
			}
			code = strings.Replace(query.SQL, "*", strings.Join(cols, ", "), 1)
		}

		gq := GoQuery{
			Cmd:          query.Cmd,
			ConstantName: lowerTitle(query.Name),
			FieldName:    lowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			SQL:          code,
		}

		if len(query.Params) == 1 {
			p := query.Params[0]
			gq.Arg = GoQueryValue{
				Name: paramName(p),
				typ:  r.goType(p.Column),
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
					Name: r.structName(paramName(p)),
					Type: r.goType(p.Column),
					Tags: map[string]string{
						"json": p.Column.Name,
					},
				})
			}
			gq.Arg = val
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
`

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
}

func lowerTitle(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func Generate(r *Result, settings GenerateSettings) (string, error) {
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
		GoQueries:           r.GoQueries(),
		Package:             pkg,
		Enums:               r.Enums(),
		Structs:             r.Structs(),
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
