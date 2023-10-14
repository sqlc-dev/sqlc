package config

import "fmt"

func Validate(c *Config) error {
	for _, sql := range c.SQL {
		if sql.Database != nil {
			if sql.Database.URI == "" && !sql.Database.Managed {
				return fmt.Errorf("invalid config: database must be managed or have a non-empty URI")
			}
		}
	}
	return nil
}
