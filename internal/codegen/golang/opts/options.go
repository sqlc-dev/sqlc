package opts

import (
	"encoding/json"
	"fmt"
	"maps"
	"path/filepath"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type Options struct {
	EmitInterface               bool              `json:"emit_interface" yaml:"emit_interface"`
	EmitJsonTags                bool              `json:"emit_json_tags" yaml:"emit_json_tags"`
	JsonTagsIdUppercase         bool              `json:"json_tags_id_uppercase" yaml:"json_tags_id_uppercase"`
	EmitDbTags                  bool              `json:"emit_db_tags" yaml:"emit_db_tags"`
	EmitPreparedQueries         bool              `json:"emit_prepared_queries" yaml:"emit_prepared_queries"`
	EmitExactTableNames         bool              `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	EmitEmptySlices             bool              `json:"emit_empty_slices,omitempty" yaml:"emit_empty_slices"`
	EmitExportedQueries         bool              `json:"emit_exported_queries" yaml:"emit_exported_queries"`
	EmitResultStructPointers    bool              `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers    bool              `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`
	EmitMethodsWithDbArgument   bool              `json:"emit_methods_with_db_argument,omitempty" yaml:"emit_methods_with_db_argument"`
	EmitPointersForNullTypes    bool              `json:"emit_pointers_for_null_types" yaml:"emit_pointers_for_null_types"`
	EmitEnumValidMethod         bool              `json:"emit_enum_valid_method,omitempty" yaml:"emit_enum_valid_method"`
	EmitAllEnumValues           bool              `json:"emit_all_enum_values,omitempty" yaml:"emit_all_enum_values"`
	EmitSqlAsComment            bool              `json:"emit_sql_as_comment,omitempty" yaml:"emit_sql_as_comment"`
	JsonTagsCaseStyle           string            `json:"json_tags_case_style,omitempty" yaml:"json_tags_case_style"`
	Package                     string            `json:"package" yaml:"package"`
	Out                         string            `json:"out" yaml:"out"`
	Overrides                   []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename                      map[string]string `json:"rename,omitempty" yaml:"rename"`
	SqlPackage                  string            `json:"sql_package" yaml:"sql_package"`
	SqlDriver                   string            `json:"sql_driver" yaml:"sql_driver"`
	OutputBatchFileName         string            `json:"output_batch_file_name,omitempty" yaml:"output_batch_file_name"`
	OutputDbFileName            string            `json:"output_db_file_name,omitempty" yaml:"output_db_file_name"`
	OutputModelsFileName        string            `json:"output_models_file_name,omitempty" yaml:"output_models_file_name"`
	OutputQuerierFileName       string            `json:"output_querier_file_name,omitempty" yaml:"output_querier_file_name"`
	OutputCopyfromFileName      string            `json:"output_copyfrom_file_name,omitempty" yaml:"output_copyfrom_file_name"`
	OutputFilesSuffix           string            `json:"output_files_suffix,omitempty" yaml:"output_files_suffix"`
	InflectionExcludeTableNames []string          `json:"inflection_exclude_table_names,omitempty" yaml:"inflection_exclude_table_names"`
	WrapErrors                  bool              `json:"wrap_errors,omitempty" yaml:"wrap_errors"`
	QueryParameterLimit         *int32            `json:"query_parameter_limit,omitempty" yaml:"query_parameter_limit"`
	OmitSqlcVersion             bool              `json:"omit_sqlc_version,omitempty" yaml:"omit_sqlc_version"`
	OmitUnusedStructs           bool              `json:"omit_unused_structs,omitempty" yaml:"omit_unused_structs"`
	BuildTags                   string            `json:"build_tags,omitempty" yaml:"build_tags"`
	Initialisms                 *[]string         `json:"initialisms,omitempty" yaml:"initialisms"`

	InitialismsMap map[string]struct{} `json:"-" yaml:"-"`
}

type GlobalOptions struct {
	Overrides []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename    map[string]string `json:"rename,omitempty" yaml:"rename"`
}

func Parse(req *plugin.GenerateRequest) (*Options, error) {
	options, err := parseOpts(req)
	if err != nil {
		return nil, err
	}
	global, err := parseGlobalOpts(req)
	if err != nil {
		return nil, err
	}
	if len(global.Overrides) > 0 {
		options.Overrides = append(global.Overrides, options.Overrides...)
	}
	if len(global.Rename) > 0 {
		if options.Rename == nil {
			options.Rename = map[string]string{}
		}
		maps.Copy(options.Rename, global.Rename)
	}
	return options, nil
}

func parseOpts(req *plugin.GenerateRequest) (*Options, error) {
	var options Options
	if len(req.PluginOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.PluginOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling plugin options: %w", err)
	}

	if options.Package == "" {
		if options.Out != "" {
			options.Package = filepath.Base(options.Out)
		} else {
			return nil, fmt.Errorf("invalid options: missing package name")
		}
	}

	for i := range options.Overrides {
		if err := options.Overrides[i].parse(req); err != nil {
			return nil, err
		}
	}

	if options.SqlPackage != "" {
		if err := validatePackage(options.SqlPackage); err != nil {
			return nil, fmt.Errorf("invalid options: %s", err)
		}
	}

	if options.SqlDriver != "" {
		if err := validateDriver(options.SqlDriver); err != nil {
			return nil, fmt.Errorf("invalid options: %s", err)
		}
	}

	if options.QueryParameterLimit == nil {
		options.QueryParameterLimit = new(int32)
		*options.QueryParameterLimit = 1
	}

	if options.Initialisms == nil {
		options.Initialisms = new([]string)
		*options.Initialisms = []string{"id"}
	}

	options.InitialismsMap = map[string]struct{}{}
	for _, initial := range *options.Initialisms {
		options.InitialismsMap[initial] = struct{}{}
	}

	return &options, nil
}

func parseGlobalOpts(req *plugin.GenerateRequest) (*GlobalOptions, error) {
	var options GlobalOptions
	if len(req.GlobalOptions) == 0 {
		return &options, nil
	}
	if err := json.Unmarshal(req.GlobalOptions, &options); err != nil {
		return nil, fmt.Errorf("unmarshalling global options: %w", err)
	}
	for i := range options.Overrides {
		if err := options.Overrides[i].parse(req); err != nil {
			return nil, err
		}
	}
	return &options, nil
}

func ValidateOpts(opts *Options) error {
	if opts.EmitMethodsWithDbArgument && opts.EmitPreparedQueries {
		return fmt.Errorf("invalid options: emit_methods_with_db_argument and emit_prepared_queries options are mutually exclusive")
	}
	if *opts.QueryParameterLimit < 0 {
		return fmt.Errorf("invalid options: query parameter limit must not be negative")
	}

	return nil
}
