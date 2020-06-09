package golang

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"

	"github.com/kyleconroy/sqlc/internal/codegen"
	"github.com/kyleconroy/sqlc/internal/config"
	"github.com/kyleconroy/sqlc/internal/metadata"
)

type Generateable interface {
	Structs(settings config.CombinedSettings) []Struct
	GoQueries(settings config.CombinedSettings) []Query
	Enums(settings config.CombinedSettings) []Enum
}

func UsesType(r Generateable, typ string, settings config.CombinedSettings) bool {
	for _, strct := range r.Structs(settings) {
		for _, f := range strct.Fields {
			fType := strings.TrimPrefix(f.Type, "[]")
			if strings.HasPrefix(fType, typ) {
				return true
			}
		}
	}
	return false
}

func UsesArrays(r Generateable, settings config.CombinedSettings) bool {
	for _, strct := range r.Structs(settings) {
		for _, f := range strct.Fields {
			if strings.HasPrefix(f.Type, "[]") {
				return true
			}
		}
	}
	return false
}

type fileImports struct {
	Std []string
	Dep []string
}

func mergeImports(imps ...fileImports) [][]string {
	if len(imps) == 1 {
		return [][]string{imps[0].Std, imps[0].Dep}
	}

	var stds, pkgs []string
	seenStd := map[string]struct{}{}
	seenPkg := map[string]struct{}{}
	for i := range imps {
		for _, std := range imps[i].Std {
			if _, ok := seenStd[std]; ok {
				continue
			}
			stds = append(stds, std)
			seenStd[std] = struct{}{}
		}
		for _, pkg := range imps[i].Dep {
			if _, ok := seenPkg[pkg]; ok {
				continue
			}
			pkgs = append(pkgs, pkg)
			seenPkg[pkg] = struct{}{}
		}
	}
	return [][]string{stds, pkgs}
}

func Imports(r Generateable, settings config.CombinedSettings) func(string) [][]string {
	return func(filename string) [][]string {
		if filename == "db.go" {
			return mergeImports(dbImports(r, settings))
		}

		if filename == "models.go" {
			return mergeImports(modelImports(r, settings))
		}

		if filename == "querier.go" {
			return mergeImports(interfaceImports(r, settings))
		}

		return mergeImports(queryImports(r, settings, filename))
	}
}

func dbImports(r Generateable, settings config.CombinedSettings) fileImports {
	std := []string{"context", "database/sql"}
	if settings.Go.EmitPreparedQueries {
		std = append(std, "fmt")
	}
	return fileImports{Std: std}
}

func interfaceImports(r Generateable, settings config.CombinedSettings) fileImports {
	gq := r.GoQueries(settings)
	uses := func(name string) bool {
		for _, q := range gq {
			if !q.Ret.isEmpty() {
				if strings.HasPrefix(q.Ret.Type(), name) {
					return true
				}
			}
			if !q.Arg.isEmpty() {
				if strings.HasPrefix(q.Arg.Type(), name) {
					return true
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
	for _, q := range gq {
		if q.Cmd == metadata.CmdExecResult {
			std["database/sql"] = struct{}{}
		}
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
	if uses("net.HardwareAddr") {
		std["net"] = struct{}{}
	}

	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
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
	return fileImports{stds, pkgs}
}

func modelImports(r Generateable, settings config.CombinedSettings) fileImports {
	std := make(map[string]struct{})
	if UsesType(r, "sql.Null", settings) {
		std["database/sql"] = struct{}{}
	}
	if UsesType(r, "json.RawMessage", settings) {
		std["encoding/json"] = struct{}{}
	}
	if UsesType(r, "time.Time", settings) {
		std["time"] = struct{}{}
	}
	if UsesType(r, "net.IP", settings) {
		std["net"] = struct{}{}
	}
	if UsesType(r, "net.HardwareAddr", settings) {
		std["net"] = struct{}{}
	}
	if len(r.Enums(settings)) > 0 {
		std["fmt"] = struct{}{}
	}

	// Custom imports
	pkg := make(map[string]struct{})
	overrideTypes := map[string]string{}
	for _, o := range settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
	}

	_, overrideNullTime := overrideTypes["pq.NullTime"]
	if UsesType(r, "pq.NullTime", settings) && !overrideNullTime {
		pkg["github.com/lib/pq"] = struct{}{}
	}

	_, overrideUUID := overrideTypes["uuid.UUID"]
	if UsesType(r, "uuid.UUID", settings) && !overrideUUID {
		pkg["github.com/google/uuid"] = struct{}{}
	}

	for goType, importPath := range overrideTypes {
		if _, ok := std[importPath]; !ok && UsesType(r, goType, settings) {
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
	return fileImports{stds, pkgs}
}

func queryImports(r Generateable, settings config.CombinedSettings, filename string) fileImports {
	// for _, strct := range r.Structs() {
	// 	for _, f := range strct.Fields {
	// 		if strings.HasPrefix(f.Type, "[]") {
	// 			return true
	// 		}
	// 	}
	// }
	var gq []Query
	for _, query := range r.GoQueries(settings) {
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
	for _, q := range gq {
		if q.Cmd == metadata.CmdExecResult {
			std["database/sql"] = struct{}{}
		}
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
	for _, o := range settings.Overrides {
		if o.GoBasicType {
			continue
		}
		overrideTypes[o.GoTypeName] = o.GoPackage
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
	return fileImports{stds, pkgs}
}

var templateSet = `
{{define "dbFile"}}// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "dbCode" . }}
{{end}}

{{define "dbCode"}}
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
{{end}}

{{define "interfaceFile"}}// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "interfaceCode" . }}
{{end}}

{{define "interfaceCode"}}
type Querier interface {
	{{- range .GoQueries}}
	{{- if eq .Cmd ":one"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ({{.Ret.Type}}, error)
	{{- end}}
	{{- if eq .Cmd ":many"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) ([]{{.Ret.Type}}, error)
	{{- end}}
	{{- if eq .Cmd ":exec"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) error
	{{- end}}
	{{- if eq .Cmd ":execrows"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (int64, error)
	{{- end}}
	{{- if eq .Cmd ":execresult"}}
	{{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (sql.Result, error)
	{{- end}}
	{{- end}}
}

var _ Querier = (*Queries)(nil)
{{end}}

{{define "modelsFile"}}// Code generated by sqlc. DO NOT EDIT.

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "modelsCode" . }}
{{end}}

{{define "modelsCode"}}
{{range .Enums}}
{{if .Comment}}{{comment .Comment}}{{end}}
type {{.Name}} string

const (
	{{- range .Constants}}
	{{.Name}} {{.Type}} = "{{.Value}}"
	{{- end}}
)

func (e *{{.Name}}) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = {{.Name}}(s)
	case string:
		*e = {{.Name}}(s)
	default:
		return fmt.Errorf("unsupported scan type for {{.Name}}: %T", src)
	}
	return nil
}
{{end}}

{{range .Structs}}
{{if .Comment}}{{comment .Comment}}{{end}}
type {{.Name}} struct { {{- range .Fields}}
  {{- if .Comment}}
  {{comment .Comment}}{{else}}
  {{- end}}
  {{.Name}} {{.Type}} {{if $.EmitJSONTags}}{{$.Q}}{{.Tag}}{{$.Q}}{{end}}
  {{- end}}
}
{{end}}
{{end}}

{{define "queryFile"}}// Code generated by sqlc. DO NOT EDIT.
// source: {{.SourceName}}

package {{.Package}}

import (
	{{range imports .SourceName}}
	{{range .}}"{{.}}"
	{{end}}
	{{end}}
)

{{template "queryCode" . }}
{{end}}

{{define "queryCode"}}
{{range .GoQueries}}
{{if $.OutputQuery .SourceName}}
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

{{if eq .Cmd ":execresult"}}
{{range .Comments}}//{{.}}
{{end -}}
func (q *Queries) {{.MethodName}}(ctx context.Context, {{.Arg.Pair}}) (sql.Result, error) {
  	{{- if $.EmitPreparedQueries}}
	return := q.exec(ctx, q.{{.FieldName}}, {{.ConstantName}}, {{.Arg.Params}})
  	{{- else}}
	return q.db.ExecContext(ctx, {{.ConstantName}}, {{.Arg.Params}})
  	{{- end}}
}
{{end}}
{{end}}
{{end}}
{{end}}
`

type tmplCtx struct {
	Q         string
	Package   string
	Enums     []Enum
	Structs   []Struct
	GoQueries []Query
	Settings  config.Config

	// TODO: Race conditions
	SourceName string

	EmitJSONTags        bool
	EmitPreparedQueries bool
	EmitInterface       bool
}

func (t *tmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

func Generate(r Generateable, settings config.CombinedSettings) (map[string]string, error) {
	funcMap := template.FuncMap{
		"lowerTitle": codegen.LowerTitle,
		"comment":    codegen.DoubleSlashComment,
		"imports":    Imports(r, settings),
	}

	tmpl := template.Must(template.New("table").Funcs(funcMap).Parse(templateSet))

	golang := settings.Go
	tctx := tmplCtx{
		Settings:            settings.Global,
		EmitInterface:       golang.EmitInterface,
		EmitJSONTags:        golang.EmitJSONTags,
		EmitPreparedQueries: golang.EmitPreparedQueries,
		Q:                   "`",
		Package:             golang.Package,
		GoQueries:           r.GoQueries(settings),
		Enums:               r.Enums(settings),
		Structs:             r.Structs(settings),
	}

	output := map[string]string{}

	execute := func(name, templateName string) error {
		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		err := tmpl.ExecuteTemplate(w, templateName, &tctx)
		w.Flush()
		if err != nil {
			return err
		}
		code, err := format.Source(b.Bytes())
		if err != nil {
			fmt.Println(b.String())
			return fmt.Errorf("source error: %w", err)
		}
		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}
		output[name] = string(code)
		return nil
	}

	if err := execute("db.go", "dbFile"); err != nil {
		return nil, err
	}
	if err := execute("models.go", "modelsFile"); err != nil {
		return nil, err
	}
	if golang.EmitInterface {
		if err := execute("querier.go", "interfaceFile"); err != nil {
			return nil, err
		}
	}

	files := map[string]struct{}{}
	for _, gq := range r.GoQueries(settings) {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}
	}
	return output, nil
}
