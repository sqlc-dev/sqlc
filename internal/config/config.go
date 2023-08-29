package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

type versionSetting struct {
	Number string `json:"version" yaml:"version"`
}

type Engine string

type Paths []string

func (p *Paths) UnmarshalJSON(data []byte) error {
	if string(data[0]) == `[` {
		var out []string
		if err := json.Unmarshal(data, &out); err != nil {
			return nil
		}
		*p = Paths(out)
		return nil
	}
	var out string
	if err := json.Unmarshal(data, &out); err != nil {
		return nil
	}
	*p = Paths([]string{out})
	return nil
}

func (p *Paths) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out := []string{}
	if sliceErr := unmarshal(&out); sliceErr != nil {
		var ele string
		if strErr := unmarshal(&ele); strErr != nil {
			return strErr
		}
		out = []string{ele}
	}

	*p = Paths(out)
	return nil
}

const (
	EngineMySQL      Engine = "mysql"
	EnginePostgreSQL Engine = "postgresql"
	EngineSQLite     Engine = "sqlite"
)

type Config struct {
	Version string   `json:"version" yaml:"version"`
	Project Project  `json:"project" yaml:"project"`
	Cloud   Cloud    `json:"cloud" yaml:"cloud"`
	SQL     []SQL    `json:"sql" yaml:"sql"`
	Gen     Gen      `json:"overrides,omitempty" yaml:"overrides"`
	Plugins []Plugin `json:"plugins" yaml:"plugins"`
	Rules   []Rule   `json:"rules" yaml:"rules"`
}

type Project struct {
	ID string `json:"id" yaml:"id"`
}

type Database struct {
	URI string `json:"uri" yaml:"uri"`
}

type Cloud struct {
	Organization string `json:"organization" yaml:"organization"`
	Project      string `json:"project" yaml:"project"`
	Hostname     string `json:"hostname" yaml:"hostname"`
}

type Plugin struct {
	Name    string   `json:"name" yaml:"name"`
	Env     []string `json:"env" yaml:"env"`
	Process *struct {
		Cmd string `json:"cmd" yaml:"cmd"`
	} `json:"process" yaml:"process"`
	WASM *struct {
		URL    string `json:"url" yaml:"url"`
		SHA256 string `json:"sha256" yaml:"sha256"`
	} `json:"wasm" yaml:"wasm"`
}

type Rule struct {
	Name string `json:"name" yaml:"name"`
	Rule string `json:"rule" yaml:"rule"`
	Msg  string `json:"message" yaml:"message"`
}

type Gen struct {
	Go *GenGo `json:"go,omitempty" yaml:"go"`
}

type GenGo struct {
	Overrides []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename    map[string]string `json:"rename,omitempty" yaml:"rename"`
}

type SQL struct {
	Engine               Engine    `json:"engine,omitempty" yaml:"engine"`
	Schema               Paths     `json:"schema" yaml:"schema"`
	Queries              Paths     `json:"queries" yaml:"queries"`
	Database             *Database `json:"database" yaml:"database"`
	StrictFunctionChecks bool      `json:"strict_function_checks" yaml:"strict_function_checks"`
	StrictOrderBy        *bool     `json:"strict_order_by" yaml:"strict_order_by"`
	Gen                  SQLGen    `json:"gen" yaml:"gen"`
	Codegen              []Codegen `json:"codegen" yaml:"codegen"`
	Rules                []string  `json:"rules" yaml:"rules"`
}

// TODO: Figure out a better name for this
type Codegen struct {
	Out     string    `json:"out" yaml:"out"`
	Plugin  string    `json:"plugin" yaml:"plugin"`
	Options yaml.Node `json:"options" yaml:"options"`
}

type SQLGen struct {
	Go   *SQLGo   `json:"go,omitempty" yaml:"go"`
	JSON *SQLJSON `json:"json,omitempty" yaml:"json"`
}

type SQLGo struct {
	EmitInterface               bool              `json:"emit_interface" yaml:"emit_interface"`
	EmitJSONTags                bool              `json:"emit_json_tags" yaml:"emit_json_tags"`
	JsonTagsIDUppercase         bool              `json:"json_tags_id_uppercase" yaml:"json_tags_id_uppercase"`
	EmitDBTags                  bool              `json:"emit_db_tags" yaml:"emit_db_tags"`
	EmitPreparedQueries         bool              `json:"emit_prepared_queries" yaml:"emit_prepared_queries"`
	EmitExactTableNames         bool              `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	EmitEmptySlices             bool              `json:"emit_empty_slices,omitempty" yaml:"emit_empty_slices"`
	EmitExportedQueries         bool              `json:"emit_exported_queries" yaml:"emit_exported_queries"`
	EmitResultStructPointers    bool              `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers    bool              `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`
	EmitMethodsWithDBArgument   bool              `json:"emit_methods_with_db_argument,omitempty" yaml:"emit_methods_with_db_argument"`
	EmitPointersForNullTypes    bool              `json:"emit_pointers_for_null_types" yaml:"emit_pointers_for_null_types"`
	EmitEnumValidMethod         bool              `json:"emit_enum_valid_method,omitempty" yaml:"emit_enum_valid_method"`
	EmitAllEnumValues           bool              `json:"emit_all_enum_values,omitempty" yaml:"emit_all_enum_values"`
	JSONTagsCaseStyle           string            `json:"json_tags_case_style,omitempty" yaml:"json_tags_case_style"`
	Package                     string            `json:"package" yaml:"package"`
	Out                         string            `json:"out" yaml:"out"`
	Overrides                   []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename                      map[string]string `json:"rename,omitempty" yaml:"rename"`
	SQLPackage                  string            `json:"sql_package" yaml:"sql_package"`
	SQLDriver                   string            `json:"sql_driver" yaml:"sql_driver"`
	OutputBatchFileName         string            `json:"output_batch_file_name,omitempty" yaml:"output_batch_file_name"`
	OutputDBFileName            string            `json:"output_db_file_name,omitempty" yaml:"output_db_file_name"`
	OutputModelsFileName        string            `json:"output_models_file_name,omitempty" yaml:"output_models_file_name"`
	OutputQuerierFileName       string            `json:"output_querier_file_name,omitempty" yaml:"output_querier_file_name"`
	OutputFilesSuffix           string            `json:"output_files_suffix,omitempty" yaml:"output_files_suffix"`
	InflectionExcludeTableNames []string          `json:"inflection_exclude_table_names,omitempty" yaml:"inflection_exclude_table_names"`
	QueryParameterLimit         *int32            `json:"query_parameter_limit,omitempty" yaml:"query_parameter_limit"`
	OmitUnusedStructs           bool              `json:"omit_unused_structs,omitempty" yaml:"omit_unused_structs"`
}

type SQLJSON struct {
	Out      string `json:"out" yaml:"out"`
	Indent   string `json:"indent,omitempty" yaml:"indent"`
	Filename string `json:"filename,omitempty" yaml:"filename"`
}

var ErrMissingEngine = errors.New("unknown engine")
var ErrMissingVersion = errors.New("no version number")
var ErrNoOutPath = errors.New("no output path")
var ErrNoPackageName = errors.New("missing package name")
var ErrNoPackagePath = errors.New("missing package path")
var ErrNoPackages = errors.New("no packages")
var ErrNoQuerierType = errors.New("no querier emit type enabled")
var ErrUnknownEngine = errors.New("invalid engine")
var ErrUnknownVersion = errors.New("invalid version number")
var ErrInvalidQueryParameterLimit = errors.New("invalid query parameter limit")

var ErrPluginBuiltin = errors.New("a built-in plugin with that name already exists")
var ErrPluginNoName = errors.New("missing plugin name")
var ErrPluginExists = errors.New("a plugin with that name already exists")
var ErrPluginNotFound = errors.New("no plugin found")
var ErrPluginNoType = errors.New("plugin: field `process` or `wasm` required")
var ErrPluginBothTypes = errors.New("plugin: both `process` and `wasm` cannot both be defined")
var ErrPluginProcessNoCmd = errors.New("plugin: missing process command")

func ParseConfig(rd io.Reader) (Config, error) {
	var buf bytes.Buffer
	var config Config
	var version versionSetting

	ver := io.TeeReader(rd, &buf)
	dec := yaml.NewDecoder(ver)
	if err := dec.Decode(&version); err != nil {
		return config, err
	}
	if version.Number == "" {
		return config, ErrMissingVersion
	}
	switch version.Number {
	case "1":
		return v1ParseConfig(&buf)
	case "2":
		return v2ParseConfig(&buf)
	default:
		return config, ErrUnknownVersion
	}
}

type CombinedSettings struct {
	Global    Config
	Package   SQL
	Go        SQLGo
	JSON      SQLJSON
	Rename    map[string]string
	Overrides []Override

	// TODO: Combine these into a more usable type
	Codegen Codegen
}

func Combine(conf Config, pkg SQL) CombinedSettings {
	cs := CombinedSettings{
		Global:  conf,
		Package: pkg,
		Rename:  map[string]string{},
	}
	if conf.Gen.Go != nil {
		for k, v := range conf.Gen.Go.Rename {
			cs.Rename[k] = v
		}
		cs.Overrides = append(cs.Overrides, conf.Gen.Go.Overrides...)
	}
	if pkg.Gen.Go != nil {
		cs.Go = *pkg.Gen.Go
		for k, v := range pkg.Gen.Go.Rename {
			cs.Rename[k] = v
		}
		cs.Overrides = append(cs.Overrides, pkg.Gen.Go.Overrides...)
	}
	if pkg.Gen.JSON != nil {
		cs.JSON = *pkg.Gen.JSON
	}
	return cs
}
