package native

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/singleflight"
)

var mysqlFlight singleflight.Group
var mysqlURI string

// StartMySQLServer installs and starts MySQL natively (without Docker).
// This is intended for CI environments like GitHub Actions where Docker may not be available.
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
		// MySQL is running without password, set one
		if err := setMySQLPassword(ctx); err != nil {
			slog.Debug("native/mysql", "set-password-error", err)
		}
		// Try again with password
		if err := waitForMySQL(ctx, uri, 1*time.Second); err == nil {
			return uri, nil
		}
		// If password didn't work, use no password
		return uriNoPassword, nil
	}

	// Try to start existing MySQL service first (might be installed but not running)
	if _, err := exec.LookPath("mysqld"); err == nil {
		slog.Info("native/mysql", "status", "starting existing service")
		if err := startMySQLService(); err == nil {
			// Wait for MySQL to be ready
			waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			// Try without password first
			if err := waitForMySQL(waitCtx, uriNoPassword, 30*time.Second); err == nil {
				if err := setMySQLPassword(ctx); err != nil {
					slog.Debug("native/mysql", "set-password-error", err)
					return uriNoPassword, nil
				}
				return uri, nil
			}

			// Try with password
			if err := waitForMySQL(waitCtx, uri, 5*time.Second); err == nil {
				return uri, nil
			}
		}
	}

	// Install MySQL if needed
	if _, err := exec.LookPath("mysql"); err != nil {
		slog.Info("native/mysql", "status", "installing")

		// Pre-configure MySQL root password
		setSelectionsCmd := exec.Command("sudo", "bash", "-c",
			`echo "mysql-server mysql-server/root_password password mysecretpassword" | sudo debconf-set-selections && `+
				`echo "mysql-server mysql-server/root_password_again password mysecretpassword" | sudo debconf-set-selections`)
		setSelectionsCmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
		if output, err := setSelectionsCmd.CombinedOutput(); err != nil {
			slog.Debug("native/mysql", "debconf", string(output))
		}

		// Try to install MySQL server
		cmd := exec.Command("sudo", "apt-get", "install", "-y", "-qq", "mysql-server")
		cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
		if output, err := cmd.CombinedOutput(); err != nil {
			// If apt-get fails (no network), return error
			return "", fmt.Errorf("apt-get install mysql-server failed (network may be unavailable): %w\n%s", err, output)
		}
	}

	// Start MySQL service
	slog.Info("native/mysql", "status", "starting service")
	if err := startMySQLService(); err != nil {
		return "", fmt.Errorf("failed to start MySQL: %w", err)
	}

	// Wait for MySQL to be ready with no password first
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Try without password first (fresh installation)
	if err := waitForMySQL(waitCtx, uriNoPassword, 30*time.Second); err == nil {
		// Set the password
		if err := setMySQLPassword(ctx); err != nil {
			slog.Debug("native/mysql", "set-password-error", err)
			// Return without password
			return uriNoPassword, nil
		}
		return uri, nil
	}

	// Try with password
	if err := waitForMySQL(waitCtx, uri, 5*time.Second); err != nil {
		return "", fmt.Errorf("timeout waiting for MySQL: %w", err)
	}

	return uri, nil
}

func startMySQLService() error {
	// Try systemctl first
	cmd := exec.Command("sudo", "systemctl", "start", "mysql")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try mysqld
	cmd = exec.Command("sudo", "systemctl", "start", "mysqld")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try service command
	cmd = exec.Command("sudo", "service", "mysql", "start")
	if err := cmd.Run(); err == nil {
		return nil
	}

	cmd = exec.Command("sudo", "service", "mysqld", "start")
	if err := cmd.Run(); err == nil {
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

	// Set root password
	_, err = db.ExecContext(ctx, "ALTER USER 'root'@'localhost' IDENTIFIED BY 'mysecretpassword';")
	if err != nil {
		// Try older MySQL syntax
		_, err = db.ExecContext(ctx, "SET PASSWORD FOR 'root'@'localhost' = PASSWORD('mysecretpassword');")
		if err != nil {
			return fmt.Errorf("could not set MySQL password: %w", err)
		}
	}

	// Flush privileges
	_, _ = db.ExecContext(ctx, "FLUSH PRIVILEGES;")

	// Create dinotest database
	_, _ = db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS dinotest;")

	return nil
}

func waitForMySQL(ctx context.Context, uri string, timeout time.Duration) error {
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
				return fmt.Errorf("timeout waiting for MySQL (last error: %v)", lastErr)
			}
			db, err := sql.Open("mysql", uri)
			if err != nil {
				lastErr = err
				slog.Debug("native/mysql", "open-attempt", err)
				continue
			}
			if err := db.PingContext(ctx); err != nil {
				lastErr = err
				db.Close()
				continue
			}
			db.Close()
			return nil
		}
	}
}
