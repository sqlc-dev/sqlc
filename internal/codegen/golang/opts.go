package golang

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sqlc-dev/sqlc/internal/plugin"
)

type opts struct {
	EmitInterface               bool     `json:"emit_interface"`
	EmitJsonTags                bool     `json:"emit_json_tags"`
	JsonTagsIdUppercase         bool     `json:"json_tags_id_uppercase"`
	EmitDbTags                  bool     `json:"emit_db_tags"`
	EmitPreparedQueries         bool     `json:"emit_prepared_queries"`
	EmitExactTableNames         bool     `json:"emit_exact_table_names,omitempty"`
	EmitEmptySlices             bool     `json:"emit_empty_slices,omitempty"`
	EmitExportedQueries         bool     `json:"emit_exported_queries"`
	EmitResultStructPointers    bool     `json:"emit_result_struct_pointers"`
	EmitParamsStructPointers    bool     `json:"emit_params_struct_pointers"`
	EmitMethodsWithDbArgument   bool     `json:"emit_methods_with_db_argument,omitempty"`
	EmitPointersForNullTypes    bool     `json:"emit_pointers_for_null_types"`
	EmitEnumValidMethod         bool     `json:"emit_enum_valid_method,omitempty"`
	EmitAllEnumValues           bool     `json:"emit_all_enum_values,omitempty"`
	JsonTagsCaseStyle           string   `json:"json_tags_case_style,omitempty"`
	Package                     string   `json:"package"`
	Out                         string   `json:"out"`
	SqlPackage                  string   `json:"sql_package"`
	SqlDriver                   string   `json:"sql_driver"`
	OutputBatchFileName         string   `json:"output_batch_file_name,omitempty"`
	OutputDbFileName            string   `json:"output_db_file_name,omitempty"`
	OutputModelsFileName        string   `json:"output_models_file_name,omitempty"`
	OutputQuerierFileName       string   `json:"output_querier_file_name,omitempty"`
	OutputCopyfromFileName      string   `json:"output_copyfrom_file_name,omitempty"`
	OutputFilesSuffix           string   `json:"output_files_suffix,omitempty"`
	InflectionExcludeTableNames []string `json:"inflection_exclude_table_names,omitempty"`
	QueryParameterLimit         *int32   `json:"query_parameter_limit,omitempty"`
	OmitUnusedStructs           bool     `json:"omit_unused_structs,omitempty"`
	BuildTags                   string   `json:"build_tags,omitempty"`

	// Unused but included in marshaled json we receive
	Overrides json.RawMessage `json:"overrides,omitempty"`
	Rename    json.RawMessage `json:"rename,omitempty"`
}

func parseOpts(req *plugin.CodeGenRequest) (*opts, error) {
	var options *opts
	dec := json.NewDecoder(bytes.NewReader(req.PluginOptions))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&options); err != nil {
		return options, fmt.Errorf("unmarshalling options: %w", err)
	}

	if options.QueryParameterLimit == nil {
		options.QueryParameterLimit = new(int32)
		*options.QueryParameterLimit = 1
	}

	return options, nil
}

func validateOpts(opts *opts) error {
	if opts.EmitMethodsWithDbArgument && opts.EmitPreparedQueries {
		return fmt.Errorf("invalid options: emit_methods_with_db_argument and emit_prepared_queries options are mutually exclusive")
	}

	return nil
}
