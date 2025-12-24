package native

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2" // ClickHouse driver
	"golang.org/x/sync/singleflight"
)

var clickhouseFlight singleflight.Group
var clickhouseURI string

// StartClickHouseServer starts an existing ClickHouse installation natively (without Docker).
func StartClickHouseServer(ctx context.Context) (string, error) {
	if err := Supported(); err != nil {
		return "", err
	}
	if clickhouseURI != "" {
		return clickhouseURI, nil
	}
	value, err, _ := clickhouseFlight.Do("clickhouse", func() (interface{}, error) {
		uri, err := startClickHouseServer(ctx)
		if err != nil {
			return "", err
		}
		clickhouseURI = uri
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

func startClickHouseServer(ctx context.Context) (string, error) {
	// Standard URI for test ClickHouse
	uri := "clickhouse://default:@localhost:9000/default"

	// Try to connect first - it might already be running
	if err := waitForClickHouse(ctx, uri, 500*time.Millisecond); err == nil {
		slog.Info("native/clickhouse", "status", "already running")
		return uri, nil
	}

	// Check if ClickHouse is installed
	if _, err := exec.LookPath("clickhouse"); err != nil {
		return "", fmt.Errorf("ClickHouse is not installed (clickhouse not found)")
	}

	// Start ClickHouse server
	slog.Info("native/clickhouse", "status", "starting service")

	if err := startClickHouseService(); err != nil {
		return "", fmt.Errorf("failed to start ClickHouse: %w", err)
	}

	// Wait for ClickHouse to be ready
	waitCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	if err := waitForClickHouse(waitCtx, uri, 30*time.Second); err != nil {
		return "", fmt.Errorf("timeout waiting for ClickHouse: %w", err)
	}

	return uri, nil
}

func startClickHouseService() error {
	// Try systemctl first
	cmd := exec.Command("sudo", "systemctl", "start", "clickhouse-server")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try service command
	cmd = exec.Command("sudo", "service", "clickhouse-server", "start")
	if err := cmd.Run(); err == nil {
		return nil
	}

	// Try running clickhouse-server directly (for single-node dev setup)
	cmd = exec.Command("clickhouse-server", "--daemon")
	if err := cmd.Run(); err == nil {
		return nil
	}

	return fmt.Errorf("could not start ClickHouse server")
}

func waitForClickHouse(ctx context.Context, uri string, timeout time.Duration) error {
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
				return fmt.Errorf("timeout waiting for ClickHouse (last error: %v)", lastErr)
			}
			conn, err := sql.Open("clickhouse", uri)
			if err != nil {
				lastErr = err
				slog.Debug("native/clickhouse", "connect-attempt", err)
				continue
			}
			if err := conn.PingContext(ctx); err != nil {
				lastErr = err
				conn.Close()
				continue
			}
			conn.Close()
			return nil
		}
	}
}
