package native

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/singleflight"
)

var postgresFlight singleflight.Group
var postgresURI string

// StartPostgreSQLServer starts an existing PostgreSQL installation natively (without Docker).
func StartPostgreSQLServer(ctx context.Context) (string, error) {
	if err := Supported(); err != nil {
		return "", err
	}
	if postgresURI != "" {
		return postgresURI, nil
	}
	value, err, _ := postgresFlight.Do("postgresql", func() (interface{}, error) {
		uri, err := startPostgreSQLServer(ctx)
		if err != nil {
			return "", err
		}
		postgresURI = uri
		return uri, nil
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

func startPostgreSQLServer(ctx context.Context) (string, error) {
	// Standard URI for test PostgreSQL
	uri := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"

	// Try to connect first - it might already be running
	if err := waitForPostgres(ctx, uri, 500*time.Millisecond); err == nil {
		slog.Info("native/postgres", "status", "already running")
		return uri, nil
	}

	// Check if PostgreSQL is installed
	if _, err := exec.LookPath("psql"); err != nil {
		return "", fmt.Errorf("PostgreSQL is not installed (psql not found)")
	}

	// Start PostgreSQL service
	slog.Info("native/postgres", "status", "starting service")

	// Try systemctl first, fall back to pg_ctlcluster
	if err := startPostgresService(); err != nil {
		return "", fmt.Errorf("failed to start PostgreSQL: %w", err)
	}

	// Configure PostgreSQL for password authentication
	if err := configurePostgres(); err != nil {
		return "", fmt.Errorf("failed to configure PostgreSQL: %w", err)
	}

	// Wait for PostgreSQL to be ready
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := waitForPostgres(waitCtx, uri, 30*time.Second); err != nil {
		return "", fmt.Errorf("timeout waiting for PostgreSQL: %w", err)
	}

	return uri, nil
}

func startPostgresService() error {
	// Try systemctl first
	cmd := exec.Command("sudo", "systemctl", "start", "postgresql")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try service command
	cmd = exec.Command("sudo", "service", "postgresql", "start")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try pg_ctlcluster (Debian/Ubuntu specific)
	// Find the installed PostgreSQL version
	output, err := exec.Command("ls", "/etc/postgresql/").CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not find PostgreSQL version: %w", err)
	}

	versions := strings.Fields(string(output))
	if len(versions) == 0 {
		return fmt.Errorf("no PostgreSQL version found in /etc/postgresql/")
	}

	version := versions[0]
	cmd = exec.Command("sudo", "pg_ctlcluster", version, "main", "start")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("pg_ctlcluster start failed: %w\n%s", err, output)
	}

	return nil
}

func configurePostgres() error {
	// Set password for postgres user using sudo -u postgres
	cmd := exec.Command("sudo", "-u", "postgres", "psql", "-c", "ALTER USER postgres PASSWORD 'postgres';")
	if output, err := cmd.CombinedOutput(); err != nil {
		// This might fail if password is already set, which is fine
		slog.Debug("native/postgres", "set-password", string(output))
	}

	// Update pg_hba.conf to allow password authentication
	// First, find the pg_hba.conf file
	output, err := exec.Command("sudo", "-u", "postgres", "psql", "-t", "-c", "SHOW hba_file;").CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not find hba_file: %w", err)
	}

	hbaFile := strings.TrimSpace(string(output))
	if hbaFile == "" {
		return fmt.Errorf("empty hba_file path")
	}

	// Check if we need to update pg_hba.conf
	catOutput, err := exec.Command("sudo", "cat", hbaFile).CombinedOutput()
	if err != nil {
		return fmt.Errorf("could not read %s: %w", hbaFile, err)
	}

	// If md5 or scram-sha-256 auth is not configured for local connections, add it
	content := string(catOutput)
	if !strings.Contains(content, "host    all             all             127.0.0.1/32            md5") &&
		!strings.Contains(content, "host    all             all             127.0.0.1/32            scram-sha-256") {

		// Prepend a rule for localhost password authentication
		newRule := "host    all             all             127.0.0.1/32            md5\n"

		// Use sed to add the rule at the beginning (after comments)
		cmd := exec.Command("sudo", "bash", "-c",
			fmt.Sprintf(`echo '%s' | cat - %s > /tmp/pg_hba.conf.new && sudo mv /tmp/pg_hba.conf.new %s`,
				newRule, hbaFile, hbaFile))
		if output, err := cmd.CombinedOutput(); err != nil {
			slog.Debug("native/postgres", "update-hba-error", string(output))
		}

		// Reload PostgreSQL to apply changes
		if err := reloadPostgres(); err != nil {
			slog.Debug("native/postgres", "reload-error", err)
		}
	}

	return nil
}

func reloadPostgres() error {
	// Try systemctl reload
	cmd := exec.Command("sudo", "systemctl", "reload", "postgresql")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try service reload
	cmd = exec.Command("sudo", "service", "postgresql", "reload")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try pg_ctlcluster reload
	output, _ := exec.Command("ls", "/etc/postgresql/").CombinedOutput()
	versions := strings.Fields(string(output))
	if len(versions) > 0 {
		cmd = exec.Command("sudo", "pg_ctlcluster", versions[0], "main", "reload")
		return cmd.Run()
	}

	return fmt.Errorf("could not reload PostgreSQL")
}

func waitForPostgres(ctx context.Context, uri string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w (last error: %v)", ctx.Err(), lastErr)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for PostgreSQL (last error: %v)", lastErr)
			}
			conn, err := pgx.Connect(ctx, uri)
			if err != nil {
				lastErr = err
				slog.Debug("native/postgres", "connect-attempt", err)
				continue
			}
			if err := conn.Ping(ctx); err != nil {
				lastErr = err
				conn.Close(ctx)
				continue
			}
			conn.Close(ctx)
			return nil
		}
	}
}
