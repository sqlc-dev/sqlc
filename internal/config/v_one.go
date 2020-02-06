package config

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
)

type v1GenerateSettings struct {
	Version   string              `json:"version"`
	Packages  []v1PackageSettings `json:"packages"`
	Overrides []Override          `json:"overrides,omitempty"`
	Rename    map[string]string   `json:"rename,omitempty"`
}

type v1PackageSettings struct {
	Name                string     `json:"name"`
	Engine              Engine     `json:"engine,omitempty"`
	Path                string     `json:"path"`
	Schema              string     `json:"schema"`
	Queries             string     `json:"queries"`
	EmitInterface       bool       `json:"emit_interface"`
	EmitJSONTags        bool       `json:"emit_json_tags"`
	EmitPreparedQueries bool       `json:"emit_prepared_queries"`
	Overrides           []Override `json:"overrides"`
}

func v1ParseConfig(rd io.Reader) (Config, error) {
	dec := json.NewDecoder(rd)
	dec.DisallowUnknownFields()
	var settings v1GenerateSettings
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

func (c *v1GenerateSettings) ValidateGlobalOverrides() error {
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

func (c *v1GenerateSettings) Translate() Config {
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
					EmitInterface:       pkg.EmitInterface,
					EmitJSONTags:        pkg.EmitJSONTags,
					EmitPreparedQueries: pkg.EmitPreparedQueries,
					Package:             pkg.Name,
					Out:                 pkg.Path,
					Overrides:           pkg.Overrides,
				},
			},
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
