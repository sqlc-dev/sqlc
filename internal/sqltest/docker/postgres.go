package docker

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

var postgresHost string

const postgresImageName = "sqlc-postgres"

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

// findRepoRoot walks up from the current directory to find the directory
// containing go.mod, which is the repository root.
func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("could not find repo root (go.mod)")
		}
		dir = parent
	}
}

// buildPostgresImage builds the fast-startup PostgreSQL image from
// Dockerfile.postgres. The Dockerfile requires no build context, so we
// pipe it to `docker build -` to avoid sending the repo tree to the daemon.
func buildPostgresImage() error {
	root, err := findRepoRoot()
	if err != nil {
		return err
	}
	content, err := os.ReadFile(filepath.Join(root, "Dockerfile.postgres"))
	if err != nil {
		return fmt.Errorf("read Dockerfile.postgres: %w", err)
	}
	cmd := exec.Command("docker", "build", "-t", postgresImageName, "-")
	cmd.Stdin = bytes.NewReader(content)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker build sqlc-postgres: %w\n%s", err, output)
	}
	return nil
}

// postgresImageExists checks whether the sqlc-postgres image is already built.
func postgresImageExists() bool {
	cmd := exec.Command("docker", "image", "inspect", postgresImageName)
	return cmd.Run() == nil
}

func startPostgreSQLServer(c context.Context) (string, error) {
	// Build the fast-startup image if it doesn't already exist.
	if !postgresImageExists() {
		if err := buildPostgresImage(); err != nil {
			return "", err
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
		// The sqlc-postgres image is pre-initialized and pre-configured,
		// so no environment variables or extra flags are needed.
		cmd := exec.Command("docker", "run",
			"--name", "sqlc_sqltest_docker_postgres",
			"-p", "5432:5432",
			"-d",
			postgresImageName,
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
