package config

import (
	"fmt"
	"io"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

func v2ParseConfig(rd io.Reader) (Config, error) {
	dec := yaml.NewDecoder(rd)
	dec.KnownFields(true)
	var conf Config
	if err := dec.Decode(&conf); err != nil {
		return conf, err
	}
	if conf.Version == "" {
		return conf, ErrMissingVersion
	}
	if conf.Version != "2" {
		return conf, ErrUnknownVersion
	}
	if len(conf.SQL) == 0 {
		return conf, ErrNoPackages
	}
	if err := conf.validateGlobalOverrides(); err != nil {
		return conf, err
	}
	if conf.Gen.Go != nil {
		for i := range conf.Gen.Go.Overrides {
			if err := conf.Gen.Go.Overrides[i].Parse(); err != nil {
				return conf, err
			}
		}
	}
	for _, sql := range conf.SQL {
		if sql.Engine == "" {
			return conf, ErrMissingEngine
		}
		if sql.Gen.Go != nil {
			if sql.Gen.Go.Out == "" {
				return conf, ErrNoPackagePath
			}
			if sql.Gen.Go.Package == "" {
				sql.Gen.Go.Package = filepath.Base(sql.Gen.Go.Out)
			}
			for i := range sql.Gen.Go.Overrides {
				if err := sql.Gen.Go.Overrides[i].Parse(); err != nil {
					return conf, err
				}
			}
		}
		if sql.Gen.Kotlin != nil {
			if sql.Gen.Kotlin.Out == "" {
				return conf, ErrNoOutPath
			}
			if sql.Gen.Kotlin.Package == "" {
				return conf, ErrNoPackageName
			}
		}
		if sql.Gen.Python != nil {
			if sql.Gen.Python.QueryParameterLimit != nil {
				if *sql.Gen.Python.QueryParameterLimit == 0 || *sql.Gen.Python.QueryParameterLimit < -1 {
					return conf, ErrInvalidQueryParameterLimit
				}
			} else {
				sql.Gen.Python.QueryParameterLimit = &defaultQueryParameterLimit
			}
			if sql.Gen.Python.Out == "" {
				return conf, ErrNoOutPath
			}
			if sql.Gen.Python.Package == "" {
				return conf, ErrNoPackageName
			}
			if !sql.Gen.Python.EmitSyncQuerier && !sql.Gen.Python.EmitAsyncQuerier {
				return conf, ErrNoQuerierType
			}
			for i := range sql.Gen.Python.Overrides {
				if err := sql.Gen.Python.Overrides[i].Parse(); err != nil {
					return conf, err
				}
			}

		}
	}
	return conf, nil
}

func (c *Config) validateGlobalOverrides() error {
	engines := map[Engine]struct{}{}
	for _, pkg := range c.SQL {
		if _, ok := engines[pkg.Engine]; !ok {
			engines[pkg.Engine] = struct{}{}
		}
	}
	if c.Gen.Go == nil {
		return nil
	}
	usesMultipleEngines := len(engines) > 1
	for _, oride := range c.Gen.Go.Overrides {
		if usesMultipleEngines && oride.Engine == "" {
			return fmt.Errorf(`the "engine" field is required for global type overrides because your configuration uses multiple database engines`)
		}
	}
	return nil
}
