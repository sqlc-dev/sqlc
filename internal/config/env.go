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

	serverUri := os.Getenv("SQLC_SERVER_URI")
	if serverUri != "" && len(c.Servers) != 1 {
		return fmt.Errorf("$SQLC_SERVER_URI may only be used when there is exactly one server in config file")
	} else if serverUri != "" {
		c.Servers[0].URI = serverUri
	}

	return nil
}
