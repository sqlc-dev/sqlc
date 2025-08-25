package docker

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

var postgresSync sync.Once
var postgresHost string

func StartPostgreSQLServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}

	{
		_, err := exec.Command("docker", "pull", "postgres:16").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: postgres:16 %w", err)
		}
	}

	var syncErr error
	postgresSync.Do(func() {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

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
		if err != nil {
			syncErr = err
			return
		}

		// Create a ticker that fires every 10ms
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		uri := "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"

		for {
			select {
			case <-ctx.Done():
				syncErr = fmt.Errorf("timeout reached: %w", ctx.Err())
				return

			case <-ticker.C:
				// Run your function here
				conn, err := pgx.Connect(ctx, uri)
				if err != nil {
					slog.Debug("sqltest", "connect", err)
					continue
				}
				if err := conn.Ping(ctx); err != nil {
					slog.Error("sqltest", "ping", err)
					continue
				}
				postgresHost = uri
				return
			}
		}
	})

	if syncErr != nil {
		return "", syncErr
	}

	if postgresHost == "" {
		return "", fmt.Errorf("postgres server setup failed")
	}

	return postgresHost, nil
}
