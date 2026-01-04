package native

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/singleflight"
)

var mysqlFlight singleflight.Group
var mysqlURI string

// StartMySQLServer starts an existing MySQL installation natively (without Docker).
func StartMySQLServer(ctx context.Context) (string, error) {
	if err := Supported(); err != nil {
		return "", err
	}
	if mysqlURI != "" {
		return mysqlURI, nil
	}
	value, err, _ := mysqlFlight.Do("mysql", func() (interface{}, error) {
		uri, err := startMySQLServer(ctx)
		if err != nil {
			return "", err
		}
		mysqlURI = uri
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

func startMySQLServer(ctx context.Context) (string, error) {
	// Standard URI for test MySQL
	uri := "root:mysecretpassword@tcp(localhost:3306)/mysql?multiStatements=true&parseTime=true"

	// Try to connect first - it might already be running
	if err := waitForMySQL(ctx, uri, 500*time.Millisecond); err == nil {
		slog.Info("native/mysql", "status", "already running")
		return uri, nil
	}

	// Also try without password (default MySQL installation)
	uriNoPassword := "root@tcp(localhost:3306)/mysql?multiStatements=true&parseTime=true"
	if err := waitForMySQL(ctx, uriNoPassword, 500*time.Millisecond); err == nil {
		slog.Info("native/mysql", "status", "already running (no password)")
		// MySQL is running without password, try to set one
		if err := setMySQLPassword(ctx); err != nil {
			slog.Debug("native/mysql", "set-password-error", err)
			// Return without password if we can't set one
			return uriNoPassword, nil
		}
		// Try again with password
		if err := waitForMySQL(ctx, uri, 1*time.Second); err == nil {
			return uri, nil
		}
		// If password didn't work, use no password
		return uriNoPassword, nil
	}

	// Try to start existing MySQL service (might be installed but not running)
	if _, err := exec.LookPath("mysqld"); err == nil {
		slog.Info("native/mysql", "status", "starting existing service")
		if err := startMySQLService(); err != nil {
			slog.Debug("native/mysql", "start-error", err)
		} else {
			// Wait for MySQL to be ready
			waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Try with password first
			if err := waitForMySQL(waitCtx, uri, 15*time.Second); err == nil {
				return uri, nil
			}

			// Try without password
			if err := waitForMySQL(waitCtx, uriNoPassword, 15*time.Second); err == nil {
				if err := setMySQLPassword(ctx); err != nil {
					slog.Debug("native/mysql", "set-password-error", err)
					return uriNoPassword, nil
				}
				if err := waitForMySQL(ctx, uri, 1*time.Second); err == nil {
					return uri, nil
				}
				return uriNoPassword, nil
			}
		}
	}

	return "", fmt.Errorf("MySQL is not installed or could not be started")
}

func startMySQLService() error {
	// Try systemctl first
	cmd := exec.Command("sudo", "systemctl", "start", "mysql")
	if err := cmd.Run(); err == nil {
		// Give MySQL time to fully initialize
		time.Sleep(2 * time.Second)
		return nil
	}

	// Try mysqld
	cmd = exec.Command("sudo", "systemctl", "start", "mysqld")
	if err := cmd.Run(); err == nil {
		time.Sleep(2 * time.Second)
		return nil
	}

	// Try service command
	cmd = exec.Command("sudo", "service", "mysql", "start")
	if err := cmd.Run(); err == nil {
		time.Sleep(2 * time.Second)
		return nil
	}

	cmd = exec.Command("sudo", "service", "mysqld", "start")
	if err := cmd.Run(); err == nil {
		time.Sleep(2 * time.Second)
		return nil
	}

	return fmt.Errorf("could not start MySQL service")
}

func setMySQLPassword(ctx context.Context) error {
	// Connect without password
	db, err := sql.Open("mysql", "root@tcp(localhost:3306)/mysql")
	if err != nil {
		return err
	}
	defer db.Close()

	// Set root password using mysql_native_password for broader compatibility
	_, err = db.ExecContext(ctx, "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'mysecretpassword';")
	if err != nil {
		// Try without specifying auth plugin
		_, err = db.ExecContext(ctx, "ALTER USER 'root'@'localhost' IDENTIFIED BY 'mysecretpassword';")
		if err != nil {
			// Try older MySQL syntax
			_, err = db.ExecContext(ctx, "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('mysecretpassword');")
			if err != nil {
				return fmt.Errorf("could not set MySQL password: %w", err)
			}
		}
	}

	// Flush privileges
	_, _ = db.ExecContext(ctx, "FLUSH PRIVILEGES;")

	return nil
}

func waitForMySQL(ctx context.Context, uri string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Make an immediate first attempt before waiting for the ticker
	if err := tryMySQLConnection(ctx, uri); err == nil {
		return nil
	}

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w (last error: %v)", ctx.Err(), lastErr)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for MySQL (last error: %v)", lastErr)
			}
			if err := tryMySQLConnection(ctx, uri); err != nil {
				lastErr = err
				continue
			}
			return nil
		}
	}
}

func tryMySQLConnection(ctx context.Context, uri string) error {
	db, err := sql.Open("mysql", uri)
	if err != nil {
		slog.Debug("native/mysql", "open-attempt", err)
		return err
	}
	defer db.Close()
	// Use a short timeout for ping to avoid hanging
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return db.PingContext(pingCtx)
}
