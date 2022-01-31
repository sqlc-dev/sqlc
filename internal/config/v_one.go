package config

import (
	"fmt"
	"io"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

type V1GenerateSettings struct {
	Version   string              `json:"version" yaml:"version"`
	Packages  []v1PackageSettings `json:"packages" yaml:"packages"`
	Overrides []Override          `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	Rename    map[string]string   `json:"rename,omitempty" yaml:"rename,omitempty"`
}

type v1PackageSettings struct {
	Name                      string     `json:"name" yaml:"name"`
	Engine                    Engine     `json:"engine,omitempty" yaml:"engine"`
	Path                      string     `json:"path" yaml:"path"`
	Schema                    Paths      `json:"schema" yaml:"schema"`
	Queries                   Paths      `json:"queries" yaml:"queries"`
	EmitInterface             bool       `json:"emit_interface" yaml:"emit_interface"`
	EmitJSONTags              bool       `json:"emit_json_tags" yaml:"emit_json_tags"`
	EmitDBTags                bool       `json:"emit_db_tags" yaml:"emit_db_tags"`
	EmitPreparedQueries       bool       `json:"emit_prepared_queries" yaml:"emit_prepared_queries"`
	EmitExactTableNames       bool       `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	EmitEmptySlices           bool       `json:"emit_empty_slices,omitempty" yaml:"emit_empty_slices"`
	EmitExportedQueries       bool       `json:"emit_exported_queries,omitempty" yaml:"emit_exported_queries"`
	EmitResultStructPointers  bool       `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers  bool       `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`
	EmitMethodsWithDBArgument bool       `json:"emit_methods_with_db_argument" yaml:"emit_methods_with_db_argument"`
	JSONTagsCaseStyle         string     `json:"json_tags_case_style,omitempty" yaml:"json_tags_case_style"`
	SQLPackage                string     `json:"sql_package" yaml:"sql_package"`
	Overrides                 []Override `json:"overrides" yaml:"overrides"`
	OutputDBFileName          string     `json:"output_db_file_name,omitempty" yaml:"output_db_file_name"`
	OutputModelsFileName      string     `json:"output_models_file_name,omitempty" yaml:"output_models_file_name"`
	OutputQuerierFileName     string     `json:"output_querier_file_name,omitempty" yaml:"output_querier_file_name"`
	OutputFilesSuffix         string     `json:"output_files_suffix,omitempty" yaml:"output_files_suffix"`
	StrictFunctionChecks      bool       `json:"strict_function_checks" yaml:"strict_function_checks"`
}

func v1ParseConfig(rd io.Reader) (Config, error) {
	dec := yaml.NewDecoder(rd)
	dec.KnownFields(true)
	var settings V1GenerateSettings
	var config Config
	if err := dec.Decode(&settings); err != nil {
		return config, err
	}
	if settings.Version == "" {
		return config, ErrMissingVersion
	}
	if settings.Version != "1" {
		return config, ErrUnknownVersion
	}
	if len(settings.Packages) == 0 {
		return config, ErrNoPackages
	}
	if err := settings.ValidateGlobalOverrides(); err != nil {
		return config, err
	}
	for i := range settings.Overrides {
		if err := settings.Overrides[i].Parse(); err != nil {
			return config, err
		}
	}
	for j := range settings.Packages {
		if settings.Packages[j].Path == "" {
			return config, ErrNoPackagePath
		}
		for i := range settings.Packages[j].Overrides {
			if err := settings.Packages[j].Overrides[i].Parse(); err != nil {
				return config, err
			}
		}
		if settings.Packages[j].Name == "" {
			settings.Packages[j].Name = filepath.Base(settings.Packages[j].Path)
		}
		if settings.Packages[j].Engine == "" {
			settings.Packages[j].Engine = EnginePostgreSQL
		}
	}
	return settings.Translate(), nil
}

func (c *V1GenerateSettings) ValidateGlobalOverrides() error {
	engines := map[Engine]struct{}{}
	for _, pkg := range c.Packages {
		if _, ok := engines[pkg.Engine]; !ok {
			engines[pkg.Engine] = struct{}{}
		}
	}

	usesMultipleEngines := len(engines) > 1
	for _, oride := range c.Overrides {
		if usesMultipleEngines && oride.Engine == "" {
			return fmt.Errorf(`the "engine" field is required for global type overrides because your configuration uses multiple database engines`)
		}
	}
	return nil
}

func (c *V1GenerateSettings) Translate() Config {
	conf := Config{
		Version: c.Version,
	}

	for _, pkg := range c.Packages {
		conf.SQL = append(conf.SQL, SQL{
			Engine:  pkg.Engine,
			Schema:  pkg.Schema,
			Queries: pkg.Queries,
			Gen: SQLGen{
				Go: &SQLGo{
					EmitInterface:             pkg.EmitInterface,
					EmitJSONTags:              pkg.EmitJSONTags,
					EmitDBTags:                pkg.EmitDBTags,
					EmitPreparedQueries:       pkg.EmitPreparedQueries,
					EmitExactTableNames:       pkg.EmitExactTableNames,
					EmitEmptySlices:           pkg.EmitEmptySlices,
					EmitExportedQueries:       pkg.EmitExportedQueries,
					EmitResultStructPointers:  pkg.EmitResultStructPointers,
					EmitParamsStructPointers:  pkg.EmitParamsStructPointers,
					EmitMethodsWithDBArgument: pkg.EmitMethodsWithDBArgument,
					Package:                   pkg.Name,
					Out:                       pkg.Path,
					SQLPackage:                pkg.SQLPackage,
					Overrides:                 pkg.Overrides,
					JSONTagsCaseStyle:         pkg.JSONTagsCaseStyle,
					OutputDBFileName:          pkg.OutputDBFileName,
					OutputModelsFileName:      pkg.OutputModelsFileName,
					OutputQuerierFileName:     pkg.OutputQuerierFileName,
					OutputFilesSuffix:         pkg.OutputFilesSuffix,
				},
			},
			StrictFunctionChecks: pkg.StrictFunctionChecks,
		})
	}

	if len(c.Overrides) > 0 || len(c.Rename) > 0 {
		conf.Gen.Go = &GenGo{
			Overrides: c.Overrides,
			Rename:    c.Rename,
		}
	}

	return conf
}
