package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"

	"gopkg.in/yaml.v3"

	golang "github.com/sqlc-dev/sqlc/internal/codegen/golang/opts"
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
	Version   string               `json:"version" yaml:"version"`
	Cloud     Cloud                `json:"cloud" yaml:"cloud"`
	Servers   []Server             `json:"servers" yaml:"servers"`
	SQL       []SQL                `json:"sql" yaml:"sql"`
	Overrides Overrides            `json:"overrides,omitempty" yaml:"overrides"`
	Plugins   []Plugin             `json:"plugins" yaml:"plugins"`
	Rules     []Rule               `json:"rules" yaml:"rules"`
	Options   map[string]yaml.Node `json:"options" yaml:"options"`
}

type Server struct {
	Name   string `json:"name,omitempty" yaml:"name"`
	Engine Engine `json:"engine,omitempty" yaml:"engine"`
	URI    string `json:"uri" yaml:"uri"`
}

type Database struct {
	URI     string `json:"uri" yaml:"uri"`
	Managed bool   `json:"managed" yaml:"managed"`
}

type Cloud struct {
	Organization string `json:"organization" yaml:"organization"`
	Project      string `json:"project" yaml:"project"`
	Hostname     string `json:"hostname" yaml:"hostname"`
	AuthToken    string `json:"-" yaml:"-"`
}

type Plugin struct {
	Name    string   `json:"name" yaml:"name"`
	Env     []string `json:"env" yaml:"env"`
	Process *struct {
		Cmd    string `json:"cmd" yaml:"cmd"`
		Format string `json:"format" yaml:"format"`
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

type Overrides struct {
	Go *golang.GlobalOptions `json:"go,omitempty" yaml:"go"`
}

type SQL struct {
	Name                 string    `json:"name" yaml:"name"`
	Engine               Engine    `json:"engine,omitempty" yaml:"engine"`
	Schema               Paths     `json:"schema" yaml:"schema"`
	Queries              Paths     `json:"queries" yaml:"queries"`
	Database             *Database `json:"database" yaml:"database"`
	StrictFunctionChecks bool      `json:"strict_function_checks" yaml:"strict_function_checks"`
	StrictOrderBy        *bool     `json:"strict_order_by" yaml:"strict_order_by"`
	Gen                  SQLGen    `json:"gen" yaml:"gen"`
	Codegen              []Codegen `json:"codegen" yaml:"codegen"`
	Rules                []string  `json:"rules" yaml:"rules"`
	Analyzer             Analyzer  `json:"analyzer" yaml:"analyzer"`
}

type Analyzer struct {
	Database *bool `json:"database" yaml:"database"`
}

// TODO: Figure out a better name for this
type Codegen struct {
	Out     string    `json:"out" yaml:"out"`
	Plugin  string    `json:"plugin" yaml:"plugin"`
	Options yaml.Node `json:"options" yaml:"options"`
}

type SQLGen struct {
	Go   *golang.Options `json:"go,omitempty" yaml:"go"`
	JSON *SQLJSON        `json:"json,omitempty" yaml:"json"`
}

type SQLJSON struct {
	Out      string `json:"out" yaml:"out"`
	Indent   string `json:"indent,omitempty" yaml:"indent"`
	Filename string `json:"filename,omitempty" yaml:"filename"`
}

var ErrMissingEngine = errors.New("unknown engine")
var ErrMissingVersion = errors.New("no version number")
var ErrNoOutPath = errors.New("no output path")
var ErrNoPackagePath = errors.New("missing package path")
var ErrNoPackages = errors.New("no packages")
var ErrNoQuerierType = errors.New("no querier emit type enabled")
var ErrUnknownEngine = errors.New("invalid engine")
var ErrUnknownVersion = errors.New("invalid version number")

var ErrPluginBuiltin = errors.New("a built-in plugin with that name already exists")
var ErrPluginNoName = errors.New("missing plugin name")
var ErrPluginExists = errors.New("a plugin with that name already exists")
var ErrPluginNotFound = errors.New("no plugin found")
var ErrPluginNoType = errors.New("plugin: field `process` or `wasm` required")
var ErrPluginBothTypes = errors.New("plugin: `process` and `wasm` cannot both be defined")
var ErrPluginProcessNoCmd = errors.New("plugin: missing process command")

var ErrInvalidDatabase = errors.New("database must be managed or have a non-empty URI")
var ErrManagedDatabaseNoProject = errors.New(`managed databases require a cloud project

If you don't have a project, you can create one from the sqlc Cloud
dashboard at https://dashboard.sqlc.dev/. If you have a project, ensure
you've set its id as the value of the "project" field within the "cloud"
section of your sqlc configuration. The id will look similar to
"01HA8TWGMYPHK0V2GGMB3R2TP9".`)
var ErrManagedDatabaseNoAuthToken = errors.New(`managed databases require an auth token

If you don't have an auth token, you can create one from the sqlc Cloud
dashboard at https://dashboard.sqlc.dev/. If you have an auth token, ensure
you've set it as the value of the SQLC_AUTH_TOKEN environment variable.`)

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
	var err error
	switch version.Number {
	case "1":
		config, err = v1ParseConfig(&buf)
		if err != nil {
			return config, err
		}
	case "2":
		config, err = v2ParseConfig(&buf)
		if err != nil {
			return config, err
		}
	default:
		return config, ErrUnknownVersion
	}
	err = config.addEnvVars()
	if err != nil {
		return config, err
	}
	return config, nil
}

type CombinedSettings struct {
	Global  Config
	Package SQL
	Go      golang.Options
	JSON    SQLJSON

	// TODO: Combine these into a more usable type
	Codegen Codegen
}

func Combine(conf Config, pkg SQL) CombinedSettings {
	cs := CombinedSettings{
		Global:  conf,
		Package: pkg,
	}
	if pkg.Gen.Go != nil {
		cs.Go = *pkg.Gen.Go
	}
	if pkg.Gen.JSON != nil {
		cs.JSON = *pkg.Gen.JSON
	}
	return cs
}
