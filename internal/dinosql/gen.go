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

	"github.com/jinzhu/inflection"
)

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
			FieldName:    lowerTitle(query.Name),
			MethodName:   query.Name,
			SQL:          code,
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
