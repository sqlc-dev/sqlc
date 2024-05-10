package config

func Validate(c *Config) error {
	for _, sql := range c.SQL {
		if sql.Database != nil {
			switch {
			case sql.Database.URI != "":
			case sql.Database.Managed:
			case sql.Database.Auto:
			default:
				return ErrInvalidDatabase
			}
		}
	}
	return nil
}
