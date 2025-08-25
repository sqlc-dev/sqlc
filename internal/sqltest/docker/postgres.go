package docker

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var postgresHost string

func StartPostgreSQLServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}
	if postgresHost != "" {
		return postgresHost, nil
	}
	value, err, _ := flight.Do("postgresql", func() (interface{}, error) {
		host, err := startPostgreSQLServer(c)
		if err != nil {
			return "", err
		}
		postgresHost = host
		return host, err
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

func startPostgreSQLServer(c context.Context) (string, error) {
	{
		_, err := exec.Command("docker", "pull", "postgres:16").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: postgres:16 %w", err)
		}
	}

	uri := "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

	var exists bool
	{
		cmd := exec.Command("docker", "container", "inspect", "sqlc_sqltest_docker_postgres")
		// This means we've already started the container
		exists = cmd.Run() == nil
	}

	if !exists {
		cmd := exec.Command("docker", "run",
			"--name", "sqlc_sqltest_docker_postgres",
			"-e", "POSTGRES_PASSWORD=mysecretpassword",
			"-e", "POSTGRES_USER=postgres",
			"-p", "5432:5432",
			"-d",
			"postgres:16",
			"-c", "max_connections=200",
		)

		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		msg := `Conflict. The container name "/sqlc_sqltest_docker_postgres" is already in use by container`
		if !strings.Contains(string(output), msg) && err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(c, 5*time.Second)
	defer cancel()

	// Create a ticker that fires every 10ms
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout reached: %w", ctx.Err())

		case <-ticker.C:
			// Run your function here
			conn, err := pgx.Connect(ctx, uri)
			if err != nil {
				slog.Debug("sqltest", "connect", err)
				continue
			}
			defer conn.Close(ctx)
			if err := conn.Ping(ctx); err != nil {
				slog.Error("sqltest", "ping", err)
				continue
			}
			return uri, nil
		}
	}
}
