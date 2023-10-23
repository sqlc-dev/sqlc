package config

func Validate(c *Config) error {
	for _, sql := range c.SQL {
		if sql.Database != nil {
			if sql.Database.URI == "" && !sql.Database.Managed {
				return ErrInvalidDatabase
			}
		}
	}
	return nil
}
