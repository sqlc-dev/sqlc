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
	for j := range conf.SQL {
		if conf.SQL[j].Engine == "" {
			return conf, ErrMissingEngine
		}
		if conf.SQL[j].Gen.Go != nil {
			if conf.SQL[j].Gen.Go.Out == "" {
				return conf, ErrNoPackagePath
			}
			if conf.SQL[j].Gen.Go.Package == "" {
				conf.SQL[j].Gen.Go.Package = filepath.Base(conf.SQL[j].Gen.Go.Out)
			}
			for i := range conf.SQL[j].Gen.Go.Overrides {
				if err := conf.SQL[j].Gen.Go.Overrides[i].Parse(); err != nil {
					return conf, err
				}
			}
		}
		if conf.SQL[j].Gen.Kotlin != nil {
			if conf.SQL[j].Gen.Kotlin.Out == "" {
				return conf, ErrNoOutPath
			}
			if conf.SQL[j].Gen.Kotlin.Package == "" {
				return conf, ErrNoPackageName
			}
		}
		if conf.SQL[j].Gen.Python != nil {
			if conf.SQL[j].Gen.Python.Out == "" {
				return conf, ErrNoOutPath
			}
			if conf.SQL[j].Gen.Python.Package == "" {
				return conf, ErrNoPackageName
			}
			if !conf.SQL[j].Gen.Python.EmitSyncQuerier && !conf.SQL[j].Gen.Python.EmitAsyncQuerier {
				return conf, ErrNoQuerierType
			}
			for i := range conf.SQL[j].Gen.Python.Overrides {
				if err := conf.SQL[j].Gen.Python.Overrides[i].Parse(); err != nil {
					return conf, err
				}
			}
		}
		if conf.SQL[j].Gen.JSON != nil {
			if conf.SQL[j].Gen.JSON.Out == "" {
				return conf, ErrNoOutPath
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
