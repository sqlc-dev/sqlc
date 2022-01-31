package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kyleconroy/sqlc/internal/sql/ast"

	yaml "gopkg.in/yaml.v3"
)

const errMessageNoVersion = `The configuration file must have a version number.
Set the version to 1 at the top of sqlc.json:

{
  "version": "1"
  ...
}
`

const errMessageUnknownVersion = `The configuration file has an invalid version number.
The only supported version is "1".
`

const errMessageNoPackages = `No packages are configured`

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

	// Experimental engines
	EngineXLemon Engine = "_lemon"
)

type Config struct {
	Version string `json:"version" yaml:"version"`
	SQL     []SQL  `json:"sql" yaml:"sql"`
	Gen     Gen    `json:"overrides,omitempty" yaml:"overrides"`
}

type Gen struct {
	Go     *GenGo     `json:"go,omitempty" yaml:"go"`
	Kotlin *GenKotlin `json:"kotlin,omitempty" yaml:"kotlin"`
}

type GenGo struct {
	Overrides []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename    map[string]string `json:"rename,omitempty" yaml:"rename"`
}

type GenKotlin struct {
	Rename map[string]string `json:"rename,omitempty" yaml:"rename"`
}

type SQL struct {
	Engine               Engine `json:"engine,omitempty" yaml:"engine"`
	Schema               Paths  `json:"schema" yaml:"schema"`
	Queries              Paths  `json:"queries" yaml:"queries"`
	StrictFunctionChecks bool   `json:"strict_function_checks" yaml:"strict_function_checks"`
	Gen                  SQLGen `json:"gen" yaml:"gen"`
}

type SQLGen struct {
	Go     *SQLGo     `json:"go,omitempty" yaml:"go"`
	Kotlin *SQLKotlin `json:"kotlin,omitempty" yaml:"kotlin"`
	Python *SQLPython `json:"python,omitempty" yaml:"python"`
}

type SQLGo struct {
	EmitInterface             bool              `json:"emit_interface" yaml:"emit_interface"`
	EmitJSONTags              bool              `json:"emit_json_tags" yaml:"emit_json_tags"`
	EmitDBTags                bool              `json:"emit_db_tags" yaml:"emit_db_tags"`
	EmitPreparedQueries       bool              `json:"emit_prepared_queries" yaml:"emit_prepared_queries"`
	EmitExactTableNames       bool              `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	EmitEmptySlices           bool              `json:"emit_empty_slices,omitempty" yaml:"emit_empty_slices"`
	EmitExportedQueries       bool              `json:"emit_exported_queries" yaml:"emit_exported_queries"`
	EmitResultStructPointers  bool              `json:"emit_result_struct_pointers" yaml:"emit_result_struct_pointers"`
	EmitParamsStructPointers  bool              `json:"emit_params_struct_pointers" yaml:"emit_params_struct_pointers"`
	EmitMethodsWithDBArgument bool              `json:"emit_methods_with_db_argument,omitempty" yaml:"emit_methods_with_db_argument"`
	JSONTagsCaseStyle         string            `json:"json_tags_case_style,omitempty" yaml:"json_tags_case_style"`
	Package                   string            `json:"package" yaml:"package"`
	Out                       string            `json:"out" yaml:"out"`
	Overrides                 []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename                    map[string]string `json:"rename,omitempty" yaml:"rename"`
	SQLPackage                string            `json:"sql_package" yaml:"sql_package"`
	OutputDBFileName          string            `json:"output_db_file_name,omitempty" yaml:"output_db_file_name"`
	OutputModelsFileName      string            `json:"output_models_file_name,omitempty" yaml:"output_models_file_name"`
	OutputQuerierFileName     string            `json:"output_querier_file_name,omitempty" yaml:"output_querier_file_name"`
	OutputFilesSuffix         string            `json:"output_files_suffix,omitempty" yaml:"output_files_suffix"`
}

type SQLKotlin struct {
	EmitExactTableNames bool   `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	Package             string `json:"package" yaml:"package"`
	Out                 string `json:"out" yaml:"out"`
}

type SQLPython struct {
	EmitExactTableNames bool       `json:"emit_exact_table_names" yaml:"emit_exact_table_names"`
	EmitSyncQuerier     bool       `json:"emit_sync_querier" yaml:"emit_sync_querier"`
	EmitAsyncQuerier    bool       `json:"emit_async_querier" yaml:"emit_async_querier"`
	Package             string     `json:"package" yaml:"package"`
	Out                 string     `json:"out" yaml:"out"`
	Overrides           []Override `json:"overrides,omitempty" yaml:"overrides"`
}

type Override struct {
	// name of the golang type to use, e.g. `github.com/segmentio/ksuid.KSUID`
	GoType GoType `json:"go_type" yaml:"go_type"`

	// name of the python type to use, e.g. `mymodule.TypeName`
	PythonType PythonType `json:"python_type" yaml:"python_type"`

	// fully qualified name of the Go type, e.g. `github.com/segmentio/ksuid.KSUID`
	DBType                  string `json:"db_type" yaml:"db_type"`
	Deprecated_PostgresType string `json:"postgres_type" yaml:"postgres_type"`

	// for global overrides only when two different engines are in use
	Engine Engine `json:"engine,omitempty" yaml:"engine"`

	// True if the GoType should override if the maching postgres type is nullable
	Nullable bool `json:"nullable" yaml:"nullable"`
	// Deprecated. Use the `nullable` property instead
	Deprecated_Null bool `json:"null" yaml:"null"`

	// fully qualified name of the column, e.g. `accounts.id`
	Column string `json:"column" yaml:"column"`

	ColumnName   *Match
	TableCatalog *Match
	TableSchema  *Match
	TableRel     *Match
	GoImportPath string
	GoPackage    string
	GoTypeName   string
	GoBasicType  bool
}

func (o *Override) Matches(n *ast.TableName, defaultSchema string) bool {
	if n == nil {
		return false
	}

	schema := n.Schema
	if n.Schema == "" {
		schema = defaultSchema
	}

	if o.TableCatalog != nil && !o.TableCatalog.MatchString(n.Catalog) {
		return false
	}

	if o.TableSchema == nil && schema != "" {
		return false
	}

	if o.TableSchema != nil && !o.TableSchema.MatchString(schema) {
		return false
	}

	if o.TableRel == nil && n.Name != "" {
		return false
	}

	if o.TableRel != nil && !o.TableRel.MatchString(n.Name) {
		return false
	}

	return true
}

func (o *Override) Parse() (err error) {

	// validate deprecated postgres_type field
	if o.Deprecated_PostgresType != "" {
		fmt.Fprintf(os.Stderr, "WARNING: \"postgres_type\" is deprecated. Instead, use \"db_type\" to specify a type override.\n")
		if o.DBType != "" {
			return fmt.Errorf(`Type override configurations cannot have "db_type" and "postres_type" together. Use "db_type" alone`)
		}
		o.DBType = o.Deprecated_PostgresType
	}

	// validate deprecated null field
	if o.Deprecated_Null {
		fmt.Fprintf(os.Stderr, "WARNING: \"null\" is deprecated. Instead, use the \"nullable\" field.\n")
		o.Nullable = true
	}

	// validate option combinations
	switch {
	case o.Column != "" && o.DBType != "":
		return fmt.Errorf("Override specifying both `column` (%q) and `db_type` (%q) is not valid.", o.Column, o.DBType)
	case o.Column == "" && o.DBType == "":
		return fmt.Errorf("Override must specify one of either `column` or `db_type`")
	}

	// validate Column
	if o.Column != "" {
		colParts := strings.Split(o.Column, ".")
		switch len(colParts) {
		case 2:
			if o.ColumnName, err = MatchCompile(colParts[1]); err != nil {
				return err
			}
			if o.TableRel, err = MatchCompile(colParts[0]); err != nil {
				return err
			}
			if o.TableSchema, err = MatchCompile("public"); err != nil {
				return err
			}
		case 3:
			if o.ColumnName, err = MatchCompile(colParts[2]); err != nil {
				return err
			}
			if o.TableRel, err = MatchCompile(colParts[1]); err != nil {
				return err
			}
			if o.TableSchema, err = MatchCompile(colParts[0]); err != nil {
				return err
			}
		case 4:
			if o.ColumnName, err = MatchCompile(colParts[3]); err != nil {
				return err
			}
			if o.TableRel, err = MatchCompile(colParts[2]); err != nil {
				return err
			}
			if o.TableSchema, err = MatchCompile(colParts[1]); err != nil {
				return err
			}
			if o.TableCatalog, err = MatchCompile(colParts[0]); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Override `column` specifier %q is not the proper format, expected '[catalog.][schema.]tablename.colname'", o.Column)
		}
	}

	// validate GoType
	parsed, err := o.GoType.Parse()
	if err != nil {
		return err
	}
	o.GoImportPath = parsed.ImportPath
	o.GoPackage = parsed.Package
	o.GoTypeName = parsed.TypeName
	o.GoBasicType = parsed.BasicType

	return nil
}

var ErrMissingVersion = errors.New("no version number")
var ErrUnknownVersion = errors.New("invalid version number")
var ErrMissingEngine = errors.New("unknown engine")
var ErrUnknownEngine = errors.New("invalid engine")
var ErrNoPackages = errors.New("no packages")
var ErrNoPackageName = errors.New("missing package name")
var ErrNoPackagePath = errors.New("missing package path")
var ErrNoOutPath = errors.New("no output path")
var ErrNoQuerierType = errors.New("no querier emit type enabled")

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

func Validate(c Config) error {
	for _, sql := range c.SQL {
		sqlGo := sql.Gen.Go
		if sqlGo == nil {
			continue
		}
		if sqlGo.EmitMethodsWithDBArgument && sqlGo.EmitPreparedQueries {
			return fmt.Errorf("invalid config: emit_methods_with_db_argument and emit_prepared_queries settings are mutually exclusive")
		}
	}
	return nil
}

type CombinedSettings struct {
	Global    Config
	Package   SQL
	Go        SQLGo
	Kotlin    SQLKotlin
	Python    SQLPython
	Rename    map[string]string
	Overrides []Override
}

func Combine(conf Config, pkg SQL) CombinedSettings {
	cs := CombinedSettings{
		Global:  conf,
		Package: pkg,
	}
	if conf.Gen.Go != nil {
		cs.Rename = conf.Gen.Go.Rename
		cs.Overrides = append(cs.Overrides, conf.Gen.Go.Overrides...)
	}
	if conf.Gen.Kotlin != nil {
		cs.Rename = conf.Gen.Kotlin.Rename
	}
	if pkg.Gen.Go != nil {
		cs.Go = *pkg.Gen.Go
		cs.Overrides = append(cs.Overrides, pkg.Gen.Go.Overrides...)
	}
	if pkg.Gen.Kotlin != nil {
		cs.Kotlin = *pkg.Gen.Kotlin
	}
	if pkg.Gen.Python != nil {
		cs.Python = *pkg.Gen.Python
		cs.Overrides = append(cs.Overrides, pkg.Gen.Python.Overrides...)
	}
	return cs
}
