package golang

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/sqlc-dev/sqlc/internal/codegen/sdk"
	"github.com/sqlc-dev/sqlc/internal/metadata"
	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type tmplCtx struct {
	Q           string
	Package     string
	SQLDriver   SQLDriver
	Enums       []Enum
	Structs     []Struct
	GoQueries   []Query
	SqlcVersion string

	// TODO: Race conditions
	SourceName string

	EmitJSONTags              bool
	JsonTagsIDUppercase       bool
	EmitDBTags                bool
	EmitPreparedQueries       bool
	EmitInterface             bool
	EmitEmptySlices           bool
	EmitMethodsWithDBArgument bool
	EmitEnumValidMethod       bool
	EmitAllEnumValues         bool
	UsesCopyFrom              bool
	UsesBatch                 bool
}

func (t *tmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

func (t *tmplCtx) codegenDbarg() string {
	if t.EmitMethodsWithDBArgument {
		return "db DBTX, "
	}
	return ""
}

// Called as a global method since subtemplate queryCodeStdExec does not have
// access to the toplevel tmplCtx
func (t *tmplCtx) codegenEmitPreparedQueries() bool {
	return t.EmitPreparedQueries
}

func (t *tmplCtx) codegenQueryMethod(q Query) string {
	db := "q.db"
	if t.EmitMethodsWithDBArgument {
		db = "db"
	}

	switch q.Cmd {
	case ":one":
		if t.EmitPreparedQueries {
			return "q.queryRow"
		}
		return db + ".QueryRowContext"

	case ":many":
		if t.EmitPreparedQueries {
			return "q.query"
		}
		return db + ".QueryContext"

	default:
		if t.EmitPreparedQueries {
			return "q.exec"
		}
		return db + ".ExecContext"
	}
}

func (t *tmplCtx) codegenQueryRetval(q Query) (string, error) {
	switch q.Cmd {
	case ":one":
		return "row :=", nil
	case ":many":
		return "rows, err :=", nil
	case ":exec":
		return "_, err :=", nil
	case ":execrows", ":execlastid":
		return "result, err :=", nil
	case ":execresult":
		return "return", nil
	default:
		return "", fmt.Errorf("unhandled q.Cmd case %q", q.Cmd)
	}
}

func Generate(ctx context.Context, req *plugin.CodeGenRequest) (*plugin.CodeGenResponse, error) {
	enums := buildEnums(req)
	structs := buildStructs(req)
	queries, err := buildQueries(req, structs)
	if err != nil {
		return nil, err
	}

	if req.Settings.Go.OmitUnusedStructs {
		enums, structs = filterUnusedStructs(enums, structs, queries)
	}

	return generate(req, enums, structs, queries)
}

func generate(req *plugin.CodeGenRequest, enums []Enum, structs []Struct, queries []Query) (*plugin.CodeGenResponse, error) {
	i := &importer{
		Settings: req.Settings,
		Queries:  queries,
		Enums:    enums,
		Structs:  structs,
	}
	golang := req.Settings.Go
	tctx := tmplCtx{
		EmitInterface:             golang.EmitInterface,
		EmitJSONTags:              golang.EmitJsonTags,
		JsonTagsIDUppercase:       golang.JsonTagsIdUppercase,
		EmitDBTags:                golang.EmitDbTags,
		EmitPreparedQueries:       golang.EmitPreparedQueries,
		EmitEmptySlices:           golang.EmitEmptySlices,
		EmitMethodsWithDBArgument: golang.EmitMethodsWithDbArgument,
		EmitEnumValidMethod:       golang.EmitEnumValidMethod,
		EmitAllEnumValues:         golang.EmitAllEnumValues,
		UsesCopyFrom:              usesCopyFrom(queries),
		UsesBatch:                 usesBatch(queries),
		SQLDriver:                 parseDriver(golang.SqlPackage, req.Settings.Engine),
		Q:                         "`",
		Package:                   golang.Package,
		Enums:                     enums,
		Structs:                   structs,
		SqlcVersion:               req.SqlcVersion,
	}

	if tctx.UsesCopyFrom && !tctx.SQLDriver.IsPGX() && tctx.SQLDriver != SQLDriverGoSQLDriverMySQL {
		return nil, errors.New(":copyfrom is only supported by pgx and github.com/go-sql-driver/mysql")
	}

	if tctx.UsesCopyFrom && tctx.SQLDriver == SQLDriverGoSQLDriverMySQL {
		if err := checkNoTimesForMySQLCopyFrom(queries); err != nil {
			return nil, err
		}
	}

	if tctx.UsesBatch && !tctx.SQLDriver.IsPGX() {
		return nil, errors.New(":batch* commands are only supported by pgx")
	}

	funcMap := template.FuncMap{
		"lowerTitle": sdk.LowerTitle,
		"comment":    sdk.DoubleSlashComment,
		"escape":     sdk.EscapeBacktick,
		"imports":    i.Imports,
		"hasPrefix":  strings.HasPrefix,

		// These methods are Go specific, they do not belong in the codegen package
		// (as that is language independent)
		"dbarg":               tctx.codegenDbarg,
		"emitPreparedQueries": tctx.codegenEmitPreparedQueries,
		"queryMethod":         tctx.codegenQueryMethod,
		"queryRetval":         tctx.codegenQueryRetval,
	}

	tmpl := template.Must(
		template.New("table").
			Funcs(funcMap).
			ParseFS(
				templates,
				"templates/*.tmpl",
				"templates/*/*.tmpl",
			),
	)

	output := map[string]string{}

	execute := func(name, templateName string) error {
		imports := i.Imports(name)
		replacedQueries := replaceConflictedArg(imports, queries)

		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		tctx.SourceName = name
		tctx.GoQueries = replacedQueries
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

		if templateName == "queryFile" && golang.OutputFilesSuffix != "" {
			name += golang.OutputFilesSuffix
		}

		if !strings.HasSuffix(name, ".go") {
			name += ".go"
		}
		output[name] = string(code)
		return nil
	}

	dbFileName := "db.go"
	if golang.OutputDbFileName != "" {
		dbFileName = golang.OutputDbFileName
	}
	modelsFileName := "models.go"
	if golang.OutputModelsFileName != "" {
		modelsFileName = golang.OutputModelsFileName
	}
	querierFileName := "querier.go"
	if golang.OutputQuerierFileName != "" {
		querierFileName = golang.OutputQuerierFileName
	}
	copyfromFileName := "copyfrom.go"
	// TODO(Jille): Make this configurable.

	batchFileName := "batch.go"
	if golang.OutputBatchFileName != "" {
		batchFileName = golang.OutputBatchFileName
	}

	if err := execute(dbFileName, "dbFile"); err != nil {
		return nil, err
	}
	if err := execute(modelsFileName, "modelsFile"); err != nil {
		return nil, err
	}
	if golang.EmitInterface {
		if err := execute(querierFileName, "interfaceFile"); err != nil {
			return nil, err
		}
	}
	if tctx.UsesCopyFrom {
		if err := execute(copyfromFileName, "copyfromFile"); err != nil {
			return nil, err
		}
	}
	if tctx.UsesBatch {
		if err := execute(batchFileName, "batchFile"); err != nil {
			return nil, err
		}
	}

	files := map[string]struct{}{}
	for _, gq := range queries {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}
	}
	resp := plugin.CodeGenResponse{}

	for filename, code := range output {
		resp.Files = append(resp.Files, &plugin.File{
			Name:     filename,
			Contents: []byte(code),
		})
	}

	return &resp, nil
}

func usesCopyFrom(queries []Query) bool {
	for _, q := range queries {
		if q.Cmd == metadata.CmdCopyFrom {
			return true
		}
	}
	return false
}

func usesBatch(queries []Query) bool {
	for _, q := range queries {
		for _, cmd := range []string{metadata.CmdBatchExec, metadata.CmdBatchMany, metadata.CmdBatchOne} {
			if q.Cmd == cmd {
				return true
			}
		}
	}
	return false
}

func checkNoTimesForMySQLCopyFrom(queries []Query) error {
	for _, q := range queries {
		for _, f := range q.Arg.Fields() {
			if f.Type == "time.Time" {
				return fmt.Errorf("values with a timezone are not yet supported")
			}
		}
	}
	return nil
}

func filterUnusedStructs(enums []Enum, structs []Struct, queries []Query) ([]Enum, []Struct) {
	keepTypes := make(map[string]struct{})

	for _, query := range queries {
		if !query.Arg.isEmpty() {
			keepTypes[query.Arg.Type()] = struct{}{}
			if query.Arg.IsStruct() {
				for _, field := range query.Arg.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
				}
			}
		}
		if query.hasRetType() {
			keepTypes[query.Ret.Type()] = struct{}{}
			if query.Ret.IsStruct() {
				for _, field := range query.Ret.Struct.Fields {
					keepTypes[field.Type] = struct{}{}
				}
			}
		}
	}

	keepEnums := make([]Enum, 0, len(enums))
	for _, enum := range enums {
		_, keep := keepTypes[enum.Name]
		_, keepNull := keepTypes["Null"+enum.Name]
		if keep || keepNull {
			keepEnums = append(keepEnums, enum)
		}
	}

	keepStructs := make([]Struct, 0, len(structs))
	for _, st := range structs {
		if _, ok := keepTypes[st.Name]; ok {
			keepStructs = append(keepStructs, st)
		}
	}

	return keepEnums, keepStructs
}
