package config

import (
	"fmt"
	"os"
	"strings"
)

func (c *Config) addEnvVars() error {
	authToken := os.Getenv("SQLC_AUTH_TOKEN")
	if authToken != "" && !strings.HasPrefix(authToken, "sqlc_") {
		return fmt.Errorf("$SQLC_AUTH_TOKEN doesn't start with \"sqlc_\"")
	}
	c.Cloud.AuthToken = authToken

	return nil
}
