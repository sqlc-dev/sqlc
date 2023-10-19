package config

func Validate(c *Config) error {
	for _, sql := range c.SQL {
		if sql.Database != nil {
			if sql.Database.URI == "" && !sql.Database.Managed {
				return ErrInvalidDatabase
			}
			if sql.Database.Managed {
				if c.Cloud.Project == "" {
					return ErrManagedDatabaseNoProject
				}
				if c.Cloud.AuthToken == "" {
					return ErrManagedDatabaseNoAuthToken
				}
			}
		}
	}
	return nil
}
