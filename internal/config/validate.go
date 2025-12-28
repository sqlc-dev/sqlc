package config

// builtinEngines contains the names of built-in database engines.
var builtinEngines = map[Engine]bool{
	EngineMySQL:      true,
	EnginePostgreSQL: true,
	EngineSQLite:     true,
}

// IsBuiltinEngine returns true if the engine name is a built-in engine.
func IsBuiltinEngine(name Engine) bool {
	return builtinEngines[name]
}

func Validate(c *Config) error {
	// Validate engine plugins
	engineNames := make(map[string]bool)
	for _, ep := range c.Engines {
		if ep.Name == "" {
			return ErrEnginePluginNoName
		}
		if IsBuiltinEngine(Engine(ep.Name)) {
			return ErrEnginePluginBuiltin
		}
		if engineNames[ep.Name] {
			return ErrEnginePluginExists
		}
		engineNames[ep.Name] = true

		if ep.Process == nil && ep.WASM == nil {
			return ErrEnginePluginNoType
		}
		if ep.Process != nil && ep.WASM != nil {
			return ErrEnginePluginBothTypes
		}
		if ep.Process != nil && ep.Process.Cmd == "" {
			return ErrEnginePluginProcessNoCmd
		}
		if ep.WASM != nil && ep.WASM.URL == "" {
			return ErrEnginePluginWASMNoURL
		}
	}

	for _, sql := range c.SQL {
		if sql.Database != nil {
			if sql.Database.URI == "" && !sql.Database.Managed {
				return ErrInvalidDatabase
			}
		}
	}
	return nil
}

// FindEnginePlugin finds an engine plugin by name.
func FindEnginePlugin(c *Config, name string) (*EnginePlugin, bool) {
	for i := range c.Engines {
		if c.Engines[i].Name == name {
			return &c.Engines[i], true
		}
	}
	return nil, false
}
