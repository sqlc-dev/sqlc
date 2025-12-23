package docker

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
	"time"

	_ "github.com/microsoft/go-mssqldb"
)

var mssqlHost string

func StartMSSQLServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}
	if mssqlHost != "" {
		return mssqlHost, nil
	}
	value, err, _ := flight.Do("mssql", func() (interface{}, error) {
		host, err := startMSSQLServer(c)
		if err != nil {
			return "", err
		}
		mssqlHost = host
		return host, nil
	})
	if err != nil {
		return "", err
	}
	data, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("returned value was not a string")
	}
	return data, nil
}

func startMSSQLServer(c context.Context) (string, error) {
	{
		_, err := exec.Command("docker", "pull", "mcr.microsoft.com/mssql/server:2022-latest").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: mssql/server:2022-latest %w", err)
		}
	}

	var exists bool
	{
		cmd := exec.Command("docker", "container", "inspect", "sqlc_sqltest_docker_mssql")
		// This means we've already started the container
		exists = cmd.Run() == nil
	}

	if !exists {
		cmd := exec.Command("docker", "run",
			"--name", "sqlc_sqltest_docker_mssql",
			"-e", "ACCEPT_EULA=Y",
			"-e", "MSSQL_SA_PASSWORD=MySecretPassword1!",
			"-e", "MSSQL_PID=Developer",
			"-p", "1433:1433",
			"-d",
			"mcr.microsoft.com/mssql/server:2022-latest",
		)

		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		msg := `Conflict. The container name "/sqlc_sqltest_docker_mssql" is already in use by container`
		if !strings.Contains(string(output), msg) && err != nil {
			return "", err
		}
	}

	// MSSQL takes longer to start than MySQL/PostgreSQL
	ctx, cancel := context.WithTimeout(c, 60*time.Second)
	defer cancel()

	// Create a ticker that fires every 500ms (MSSQL takes longer to start)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	uri := "sqlserver://sa:MySecretPassword1!@localhost:1433?database=master"

	db, err := sql.Open("sqlserver", uri)
	if err != nil {
		return "", fmt.Errorf("sql.Open: %w", err)
	}

	defer db.Close()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout reached: %w", ctx.Err())

		case <-ticker.C:
			// Run your function here
			if err := db.PingContext(ctx); err != nil {
				continue
			}
			return uri, nil
		}
	}
}
