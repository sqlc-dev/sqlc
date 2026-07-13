package config

func Validate(c *Config) error {
	for _, sql := range c.SQL {
		if sql.Database != nil {
			if sql.Database.URI == "" && !sql.Database.Managed && sql.Database.TestcontainersImage == "" {
				return ErrInvalidDatabase
			}
			if sql.Database.URI != "" && sql.Database.TestcontainersImage != "" {
				return ErrDatabaseConflict
			}
			if sql.Database.Managed && sql.Database.TestcontainersImage != "" {
				return ErrManagedTestcontainersConflict
			}
		}
	}
	return nil
}
