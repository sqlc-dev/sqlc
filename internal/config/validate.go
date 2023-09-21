package config

import "fmt"

func Validate(c *Config) error {
	for _, sql := range c.SQL {
		sqlGo := sql.Gen.Go
		if sqlGo == nil {
			continue
		}
		if sqlGo.EmitMethodsWithDBArgument && sqlGo.EmitPreparedQueries {
			return fmt.Errorf("invalid config: emit_methods_with_db_argument and emit_prepared_queries settings are mutually exclusive")
		}
		if sql.Database != nil {
			if sql.Database.URI == "" && !sql.Database.Managed {
				return fmt.Errorf("invalid config: database must be managed or have a non-empty URI")
			}
		}
	}
	return nil
}
