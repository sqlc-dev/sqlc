package native

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	_ "github.com/microsoft/go-mssqldb"
	"golang.org/x/sync/singleflight"
)

var mssqlFlight singleflight.Group
var mssqlURI string

// StartMSSQLServer starts an existing MSSQL Server installation natively (without Docker).
func StartMSSQLServer(ctx context.Context) (string, error) {
	if err := Supported(); err != nil {
		return "", err
	}
	if mssqlURI != "" {
		return mssqlURI, nil
	}
	value, err, _ := mssqlFlight.Do("mssql", func() (interface{}, error) {
		uri, err := startMSSQLServer(ctx)
		if err != nil {
			return "", err
		}
		mssqlURI = uri
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

func startMSSQLServer(ctx context.Context) (string, error) {
	// Standard URI for test MSSQL - matches docker-compose.yml password
	uri := "sqlserver://sa:MySecretPassword1!@localhost:1433?database=master"

	// Try to connect first - it might already be running
	if err := waitForMSSQL(ctx, uri, 500*time.Millisecond); err == nil {
		slog.Info("native/mssql", "status", "already running")
		return uri, nil
	}

	// Check if MSSQL is installed
	if _, err := exec.LookPath("sqlservr"); err != nil {
		// Also check for the mssql-conf tool
		if _, err := exec.LookPath("/opt/mssql/bin/mssql-conf"); err != nil {
			return "", fmt.Errorf("MSSQL Server is not installed")
		}
	}

	// Try to start existing MSSQL service
	slog.Info("native/mssql", "status", "starting existing service")
	if err := startMSSQLService(); err != nil {
		slog.Debug("native/mssql", "start-error", err)
	} else {
		// Wait for MSSQL to be ready
		waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		if err := waitForMSSQL(waitCtx, uri, 30*time.Second); err == nil {
			return uri, nil
		}
	}

	return "", fmt.Errorf("MSSQL Server is not installed or could not be started")
}

func startMSSQLService() error {
	// Try systemctl first
	cmd := exec.Command("sudo", "systemctl", "start", "mssql-server")
	if err := cmd.Run(); err == nil {
		// Give MSSQL time to fully initialize
		time.Sleep(3 * time.Second)
		return nil
	}

	// Try service command
	cmd = exec.Command("sudo", "service", "mssql-server", "start")
	if err := cmd.Run(); err == nil {
		time.Sleep(3 * time.Second)
		return nil
	}

	return fmt.Errorf("could not start MSSQL service")
}

func waitForMSSQL(ctx context.Context, uri string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	var lastErr error
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled: %w (last error: %v)", ctx.Err(), lastErr)
		case <-ticker.C:
			if time.Now().After(deadline) {
				return fmt.Errorf("timeout waiting for MSSQL (last error: %v)", lastErr)
			}
			db, err := sql.Open("sqlserver", uri)
			if err != nil {
				lastErr = err
				slog.Debug("native/mssql", "open-attempt", err)
				continue
			}
			// Use a short timeout for ping to avoid hanging
			pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
			err = db.PingContext(pingCtx)
			cancel()
			if err != nil {
				lastErr = err
				db.Close()
				continue
			}
			db.Close()
			return nil
		}
	}
}
