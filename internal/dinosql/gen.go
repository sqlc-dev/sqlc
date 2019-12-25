package dinosql

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"
	"unicode"

	core "github.com/kyleconroy/sqlc/internal/pg"

	"github.com/jinzhu/inflection"
)

var identPattern = regexp.MustCompile("[^a-zA-Z0-9]+")

type GoConstant struct {
	Name  string
	Type  string
	Value string
}

type GoEnum struct {
	Name      string
	Comment   string
	Constants []GoConstant
}

type GoField struct {
	Name    string
	Type    string
	Tags    map[string]string
	Comment string
}

func (gf GoField) Tag() string {
	tags := make([]string, 0, len(gf.Tags))
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
	Table   *core.FQN
	Name    string
	Fields  []GoField
	Comment string
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
		if strings.HasPrefix(v.typ, "[]") && v.typ != "[]byte" {
			out = append(out, "pq.Array("+v.Name+")")
		} else {
			out = append(out, v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
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
		if strings.HasPrefix(v.typ, "[]") && v.typ != "[]byte" {
			out = append(out, "pq.Array(&"+v.Name+")")
		} else {
			out = append(out, "&"+v.Name)
		}
	} else {
		for _, f := range v.Struct.Fields {
			if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
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
	Comments     []string
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
			fType := strings.TrimPrefix(f.Type, "[]")
			if strings.HasPrefix(fType, typ) {
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

func (r Result) Imports(settings PackageSettings) func(string) [][]string {
	return func(filename string) [][]string {
		if filename == "db.go" {
			imps := []string{"context", "database/sql"}
			if settings.EmitPreparedQueries {
				imps = append(imps, "fmt")
			}
			return [][]string{imps}
		}

		if filename == "models.go" {
			return r.ModelImports()
		}

		return r.QueryImports(filename)
	}
}

func (r Result) ModelImports() [][]string {
	std := make(map[string]struct{})
	if r.UsesType("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	if r.UsesType("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if r.UsesType("time.Time") {
		std["time"] = struct{}{}
	}
	if r.UsesType("net.IP") {
		std["net"] = struct{}{}
	}

	// Custom imports
	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range append(r.Settings.Overrides, r.packageSettings.Overrides...) {
		overrideTypes[o.goTypeName] = o.goPackage
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if r.UsesType("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}

	_, overrideUUID := overrideTypes["uuid.UUID"]
	if r.UsesType("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && r.UsesType(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return [][]string{stds, pkgs}
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
						fType := strings.TrimPrefix(f.Type, "[]")
						if strings.HasPrefix(fType, name) {
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
						fType := strings.TrimPrefix(f.Type, "[]")
						if strings.HasPrefix(fType, name) {
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
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Ret.Type(), "[]") && q.Ret.Type() != "[]byte" {
						return true
					}
				}
			}
			if !q.Arg.isEmpty() {
				if q.Arg.IsStruct() {
					for _, f := range q.Arg.Struct.Fields {
						if strings.HasPrefix(f.Type, "[]") && f.Type != "[]byte" {
							return true
						}
					}
				} else {
					if strings.HasPrefix(q.Arg.Type(), "[]") && q.Arg.Type() != "[]byte" {
						return true
					}
				}
			}
		}
		return false
	}

	std := map[string]struct{}{
		"context": struct{}{},
	}
	if uses("sql.Null") {
		std["database/sql"] = struct{}{}
	}
	if uses("json.RawMessage") {
		std["encoding/json"] = struct{}{}
	}
	if uses("time.Time") {
		std["time"] = struct{}{}
	}
	if uses("net.IP") {
		std["net"] = struct{}{}
	}

	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range append(r.Settings.Overrides, r.packageSettings.Overrides...) {
		overrideTypes[o.goTypeName] = o.goPackage
	}

	if sliceScan() {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if uses("pq.NullTime") && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}
	_, overrideUUID := overrideTypes["uuid.UUID"]
	if uses("uuid.UUID") && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	// Custom imports
	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && uses(goType) {
			pkg[importPath] = struct{}{}
		}
	}

	pkgs := make([]string, 0, len(pkg))
	for p, _ := range pkg {
		pkgs = append(pkgs, p)
	}

	stds := make([]string, 0, len(std))
	for s, _ := range std {
		stds = append(stds, s)
	}

	sort.Strings(stds)
	sort.Strings(pkgs)
	return [][]string{stds, pkgs}
}

func (r Result) Enums() []GoEnum {
	var enums []GoEnum
	for name, schema := range r.Catalog.Schemas {
		if name == "pg_catalog" {
			continue
		}
		for _, enum := range schema.Enums {
			var enumName string
			if name == "public" {
				enumName = enum.Name
			} else {
				enumName = name + "_" + enum.Name
			}
			e := GoEnum{
				Name:    r.structName(enumName),
				Comment: enum.Comment,
			}
			for _, v := range enum.Vals {
				name := ""
				id := strings.Replace(v, "-", "_", -1)
				id = strings.Replace(id, ":", "_", -1)
				id = strings.Replace(id, "/", "_", -1)
				id = identPattern.ReplaceAllString(id, "")
				for _, part := range strings.Split(id, "_") {
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
	if rename := r.Settings.Rename[name]; rename != "" {
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
			var tableName string
			if name == "public" {
				tableName = table.Name
			} else {
				tableName = name + "_" + table.Name
			}
			s := GoStruct{
				Table:   &core.FQN{Schema: name, Rel: table.Name},
				Name:    inflection.Singular(r.structName(tableName)),
				Comment: table.Comment,
			}
			for _, column := range table.Columns {
				s.Fields = append(s.Fields, GoField{
					Name:    r.structName(column.Name),
					Type:    r.goType(column),
					Tags:    map[string]string{"json:": column.Name},
					Comment: column.Comment,
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
	// package overrides have a higher precedence
	for _, oride := range append(r.Settings.Overrides, r.packageSettings.Overrides...) {
		if oride.Column != "" && oride.columnName == col.Name && oride.table == col.Table {
			return oride.goTypeName
		}
	}
	typ := r.goInnerType(col)
	if col.IsArray {
		return "[]" + typ
	}
	return typ
}

func (r Result) goInnerType(col core.Column) string {
	columnType := col.DataType
	notNull := col.NotNull || col.IsArray

	// package overrides have a higher precedence
	for _, oride := range append(r.Settings.Overrides, r.packageSettings.Overrides...) {
		if oride.PostgresType != "" && oride.PostgresType == columnType && oride.Null != notNull {
			return oride.goTypeName
		}
	}

	switch columnType {
	case "serial", "pg_catalog.serial4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigserial", "pg_catalog.serial8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallserial", "pg_catalog.serial2":
		return "int16"

	case "integer", "int", "pg_catalog.int4":
		if notNull {
			return "int32"
		}
		return "sql.NullInt32"

	case "bigint", "pg_catalog.int8":
		if notNull {
			return "int64"
		}
		return "sql.NullInt64"

	case "smallint", "pg_catalog.int2":
		return "int16"

	case "float", "double precision", "pg_catalog.float8":
		if notNull {
			return "float64"
		}
		return "sql.NullFloat64"

	case "real", "pg_catalog.float4":
		if notNull {
			return "float32"
		}
		return "sql.NullFloat64" // TODO: Change to sql.NullFloat32 after updating the go.mod file

	case "bool", "pg_catalog.bool":
		if notNull {
			return "bool"
		}
		return "sql.NullBool"

	case "jsonb":
		return "json.RawMessage"

	case "bytea", "blob", "pg_catalog.bytea":
		return "[]byte"

	case "date":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.time", "pg_catalog.timetz":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "pg_catalog.timestamp", "pg_catalog.timestamptz", "timestamptz":
		if notNull {
			return "time.Time"
		}
		return "sql.NullTime"

	case "text", "pg_catalog.varchar", "pg_catalog.bpchar", "string":
		if notNull {
			return "string"
		}
		return "sql.NullString"

	case "uuid":
		return "uuid.UUID"

	case "inet":
		return "net.IP"

	case "void":
		// A void value always returns NULL. Since there is no built-in NULL
		// value into the SQL package, we'll use sql.NullBool
		return "sql.NullBool"

	case "any":
		return "interface{}"

	default:
		for name, schema := range r.Catalog.Schemas {
			if name == "pg_catalog" {
				continue
			}
			for _, enum := range schema.Enums {
				if columnType == enum.Name {
					if name == "public" {
						return r.structName(enum.Name)
					}

					return r.structName(name + "_" + enum.Name)
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
		seen[c.Name]++
	}
	return &gs
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

func paramName(p Parameter) string {
	if p.Column.Name != "" {
		return argName(p.Column.Name)
	}
	return fmt.Sprintf("dollar_%d", p.Number)
}

func columnName(c core.Column, pos int) string {
	if c.Name != "" {
		return c.Name
	}
	return fmt.Sprintf("column_%d", pos+1)
}

func compareFQN(a *core.FQN, b *core.FQN) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Catalog == b.Catalog && a.Schema == b.Schema && a.Rel == b.Rel
}

func (r Result) GoQueries() []GoQuery {
	structs := r.Structs()

	qs := make([]GoQuery, 0, len(r.Queries))
	for _, query := range r.Queries {
		if query.Name == "" {
			continue
		}
		if query.Cmd == "" {
			continue
		}

		gq := GoQuery{
			Cmd:          query.Cmd,
			ConstantName: lowerTitle(query.Name),
			FieldName:    lowerTitle(query.Name) + "Stmt",
			MethodName:   query.Name,
			SourceName:   query.Filename,
			SQL:          query.SQL,
			Comments:     query.Comments,
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
					sameFQN := compareFQN(s.Table, &c.Table)
					if !sameName || !sameType || !sameFQN {
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

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

{{if .EmitPreparedQueries}}
func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	{{- if eq (len .GoQueries) 0 }}
	_ = err
	{{- end }}
	{{- range .GoQueries }}
	if q.{{.FieldName}}, err = db.PrepareContext(ctx, {{.ConstantName}}); err != nil {
		return nil, fmt.Errorf("error preparing query {{.MethodName}}: %w", err)
	}
	{{- end}}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	{{- range .GoQueries }}
	if q.{{.FieldName}} != nil {
		if cerr := q.{{.FieldName}}.Close(); cerr != nil {
			err = fmt.Errorf("error closing {{.FieldName}}: %w", cerr)
		}
	}
	{{- end}}
	return err
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
	db DBTX

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
{{if .Comment}}// {{.Comment}}{{end}}
type {{.Name}} string

const (
	{{- range .Constants}}
	{{.Name}} {{.Type}} = "{{.Value}}"
	{{- end}}
)

func (e *{{.Name}}) Scan(src interface{}) error {
	*e = {{.Name}}(src.([]byte))
	return nil
}
{{end}}

{{range .Structs}}
{{if .Comment}}// {{.Comment}}{{end}}
type {{.Name}} struct { {{- range .Fields}}
  {{- if .Comment}}
  // {{.Comment}}{{else}}
  {{- end}}
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
const {{.ConstantName}} = {{$.Q}}-- name: {{.MethodName}} {{.Cmd}}
{{.SQL}}
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
{{range .Comments}}//{{.}}
{{end -}}
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
{{range .Comments}}//{{.}}
{{end -}}
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
{{range .Comments}}//{{.}}
{{end -}}
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
{{range .Comments}}//{{.}}
{{end -}}
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
	r.packageSettings = settings
	funcMap := template.FuncMap{
		"lowerTitle": lowerTitle,
		"imports":    r.Imports(settings),
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

	for source := range files {
		if err := execute(source, sqlFile); err != nil {
			return nil, err
		}
	}
	return output, nil
}
