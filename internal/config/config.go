package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/types"
	"io"
	"os"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/kyleconroy/sqlc/internal/core"
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
	EngineMySQLBeta  Engine = "mysql:beta"
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
	Engine  Engine `json:"engine,omitempty" yaml:"engine"`
	Schema  Paths  `json:"schema" yaml:"schema"`
	Queries Paths  `json:"queries" yaml:"queries"`
	Gen     SQLGen `json:"gen" yaml:"gen"`
}

type SQLGen struct {
	Go     *SQLGo     `json:"go,omitempty" yaml:"go"`
	Kotlin *SQLKotlin `json:"kotlin,omitempty" yaml:"kotlin"`
}

type SQLGo struct {
	EmitInterface       bool              `json:"emit_interface" yaml:"emit_interface"`
	EmitJSONTags        bool              `json:"emit_json_tags" yaml:"emit_json_tags"`
	EmitDBTags          bool              `json:"emit_db_tags" yaml:"emit_db_tags"`
	EmitPreparedQueries bool              `json:"emit_prepared_queries" yaml:"emit_prepared_queries"`
	EmitExactTableNames bool              `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	EmitEmptySlices     bool              `json:"emit_empty_slices,omitempty" yaml:"emit_empty_slices"`
	Package             string            `json:"package" yaml:"package"`
	Out                 string            `json:"out" yaml:"out"`
	Overrides           []Override        `json:"overrides,omitempty" yaml:"overrides"`
	Rename              map[string]string `json:"rename,omitempty" yaml:"rename"`
}

type SQLKotlin struct {
	EmitExactTableNames bool   `json:"emit_exact_table_names,omitempty" yaml:"emit_exact_table_names"`
	Package             string `json:"package" yaml:"package"`
	Out                 string `json:"out" yaml:"out"`
}

type Override struct {
	// name of the golang type to use, e.g. `github.com/segmentio/ksuid.KSUID`
	GoType string `json:"go_type" yaml:"go_type"`

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

	ColumnName  string
	Table       core.FQN
	GoTypeName  string
	GoPackage   string
	GoBasicType bool
}

func (o *Override) Parse() error {

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
			o.ColumnName = colParts[1]
			o.Table = core.FQN{Schema: "public", Rel: colParts[0]}
		case 3:
			o.ColumnName = colParts[2]
			o.Table = core.FQN{Schema: colParts[0], Rel: colParts[1]}
		case 4:
			o.ColumnName = colParts[3]
			o.Table = core.FQN{Catalog: colParts[0], Schema: colParts[1], Rel: colParts[2]}
		default:
			return fmt.Errorf("Override `column` specifier %q is not the proper format, expected '[catalog.][schema.]colname.tablename'", o.Column)
		}
	}

	// validate GoType
	lastDot := strings.LastIndex(o.GoType, ".")
	lastSlash := strings.LastIndex(o.GoType, "/")
	typename := o.GoType
	if lastDot == -1 && lastSlash == -1 {
		// if the type name has no slash and no dot, validate that the type is a basic Go type
		var found bool
		for _, typ := range types.Typ {
			info := typ.Info()
			if info == 0 {
				continue
			}
			if info&types.IsUntyped != 0 {
				continue
			}
			if typename == typ.Name() {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("Package override `go_type` specifier %q is not a Go basic type e.g. 'string'", o.GoType)
		}
		o.GoBasicType = true
	} else {
		// assume the type lives in a Go package
		if lastDot == -1 {
			return fmt.Errorf("Package override `go_type` specifier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", o.GoType)
		}
		if lastSlash == -1 {
			return fmt.Errorf("Package override `go_type` specifier %q is not the proper format, expected 'package.type', e.g. 'github.com/segmentio/ksuid.KSUID'", o.GoType)
		}
		typename = o.GoType[lastSlash+1:]
		if strings.HasPrefix(typename, "go-") {
			// a package name beginning with "go-" will give syntax errors in
			// generated code. We should do the right thing and get the actual
			// import name, but in lieu of that, stripping the leading "go-" may get
			// us what we want.
			typename = typename[len("go-"):]
		}
		if strings.HasSuffix(typename, "-go") {
			typename = typename[:len(typename)-len("-go")]
		}
		o.GoPackage = o.GoType[:lastDot]
	}
	o.GoTypeName = typename
	isPointer := o.GoType[0] == '*'
	if isPointer {
		o.GoPackage = o.GoPackage[1:]
		o.GoTypeName = "*" + o.GoTypeName
	}

	return nil
}

var ErrMissingVersion = errors.New("no version number")
var ErrUnknownVersion = errors.New("invalid version number")
var ErrMissingEngine = errors.New("unknown engine")
var ErrUnknownEngine = errors.New("invalid engine")
var ErrNoPackages = errors.New("no packages")
var ErrNoPackageName = errors.New("missing package name")
var ErrNoPackagePath = errors.New("missing package path")
var ErrKotlinNoOutPath = errors.New("no output path")

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
	Kotlin    SQLKotlin
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
	return cs
}
