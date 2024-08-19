package config

import (
	"fmt"
	"io"

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
	// TODO: Store built-in plugins somewhere else
	builtins := map[string]struct{}{
		"go":   {},
		"json": {},
	}
	plugins := map[string]struct{}{}
	for i := range conf.Plugins {
		if conf.Plugins[i].Name == "" {
			return conf, ErrPluginNoName
		}
		if _, ok := builtins[conf.Plugins[i].Name]; ok {
			return conf, ErrPluginBuiltin
		}
		if _, ok := plugins[conf.Plugins[i].Name]; ok {
			return conf, ErrPluginExists
		}
		if conf.Plugins[i].Process == nil && conf.Plugins[i].WASM == nil {
			return conf, ErrPluginNoType
		}
		if conf.Plugins[i].Process != nil && conf.Plugins[i].WASM != nil {
			return conf, ErrPluginBothTypes
		}
		if conf.Plugins[i].Process != nil {
			if conf.Plugins[i].Process.Cmd == "" {
				return conf, ErrPluginProcessNoCmd
			}
		}
		plugins[conf.Plugins[i].Name] = struct{}{}
	}
	for j := range conf.SQL {
		if conf.SQL[j].Engine == "" {
			return conf, ErrMissingEngine
		}
		if conf.SQL[j].Gen.Go != nil {
			if conf.SQL[j].Gen.Go.Out == "" {
				return conf, ErrNoPackagePath
			}
		}
		if conf.SQL[j].Gen.JSON != nil {
			if conf.SQL[j].Gen.JSON.Out == "" {
				return conf, ErrNoOutPath
			}
		}
		for _, cg := range conf.SQL[j].Codegen {
			if cg.Plugin == "" {
				return conf, ErrPluginNoName
			}
			if cg.Out == "" {
				return conf, ErrNoOutPath
			}
			// TODO: Allow the use of built-in codegen from here
			if _, ok := plugins[cg.Plugin]; !ok {
				return conf, ErrPluginNotFound
			}
		}
		if conf.SQL[j].StrictOrderBy == nil {
			defaultValidate := true
			conf.SQL[j].StrictOrderBy = &defaultValidate
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
	if c.Overrides.Go == nil {
		return nil
	}
	usesMultipleEngines := len(engines) > 1
	for _, oride := range c.Overrides.Go.Overrides {
		if usesMultipleEngines && oride.Engine == "" {
			return fmt.Errorf(`the "engine" field is required for global type overrides because your configuration uses multiple database engines`)
		}
	}
	return nil
}
