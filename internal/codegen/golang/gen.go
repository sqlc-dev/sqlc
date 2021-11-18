package golang

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/kyleconroy/sqlc/internal/codegen"
	"github.com/kyleconroy/sqlc/internal/compiler"
	"github.com/kyleconroy/sqlc/internal/config"
)

type Generateable interface {
	Structs(settings config.CombinedSettings) []Struct
	GoQueries(settings config.CombinedSettings) []Query
	Enums(settings config.CombinedSettings) []Enum
}

type tmplCtx struct {
	Q          string
	Package    string
	SQLPackage SQLPackage
	Enums      []Enum
	Structs    []Struct
	GoQueries  []Query
	Settings   config.Config

	// TODO: Race conditions
	SourceName string

	EmitJSONTags              bool
	EmitDBTags                bool
	EmitPreparedQueries       bool
	EmitInterface             bool
	EmitEmptySlices           bool
	EmitMethodsWithDBArgument bool
}

func (t *tmplCtx) OutputQuery(sourceName string) bool {
	return t.SourceName == sourceName
}

func Generate(r *compiler.Result, settings config.CombinedSettings) (map[string]string, error) {
	enums := buildEnums(r, settings)
	structs := buildStructs(r, settings)
	queries := buildQueries(r, settings, structs)
	return generate(settings, enums, structs, queries)
}

func generate(settings config.CombinedSettings, enums []Enum, structs []Struct, queries []Query) (map[string]string, error) {
	i := &importer{
		Settings: settings,
		Queries:  queries,
		Enums:    enums,
		Structs:  structs,
	}

	funcMap := template.FuncMap{
		"lowerTitle": codegen.LowerTitle,
		"comment":    codegen.DoubleSlashComment,
		"escape":     codegen.EscapeBacktick,
		"imports":    i.Imports,
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

	golang := settings.Go
	tctx := tmplCtx{
		Settings:                  settings.Global,
		EmitInterface:             golang.EmitInterface,
		EmitJSONTags:              golang.EmitJSONTags,
		EmitDBTags:                golang.EmitDBTags,
		EmitPreparedQueries:       golang.EmitPreparedQueries,
		EmitEmptySlices:           golang.EmitEmptySlices,
		EmitMethodsWithDBArgument: golang.EmitMethodsWithDBArgument,
		SQLPackage:                SQLPackageFromString(golang.SQLPackage),
		Q:                         "`",
		Package:                   golang.Package,
		GoQueries:                 queries,
		Enums:                     enums,
		Structs:                   structs,
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
	if golang.OutputDBFileName != "" {
		dbFileName = golang.OutputDBFileName
	}
	modelsFileName := "models.go"
	if golang.OutputModelsFileName != "" {
		modelsFileName = golang.OutputModelsFileName
	}
	querierFileName := "querier.go"
	if golang.OutputQuerierFileName != "" {
		querierFileName = golang.OutputQuerierFileName
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

	files := map[string]struct{}{}
	for _, gq := range queries {
		files[gq.SourceName] = struct{}{}
	}

	for source := range files {
		if err := execute(source, "queryFile"); err != nil {
			return nil, err
		}
	}
	return output, nil
}
