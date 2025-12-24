package docker

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2" // ClickHouse driver
)

var clickhouseHost string

func StartClickHouseServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}
	if clickhouseHost != "" {
		return clickhouseHost, nil
	}
	value, err, _ := flight.Do("clickhouse", func() (interface{}, error) {
		host, err := startClickHouseServer(c)
		if err != nil {
			return "", err
		}
		clickhouseHost = host
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

func startClickHouseServer(c context.Context) (string, error) {
	{
		_, err := exec.Command("docker", "pull", "clickhouse/clickhouse-server:latest").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: clickhouse/clickhouse-server:latest %w", err)
		}
	}

	uri := "clickhouse://default:mysecretpassword@localhost:9000/default"

	var exists bool
	{
		cmd := exec.Command("docker", "container", "inspect", "sqlc_sqltest_docker_clickhouse")
		exists = cmd.Run() == nil
	}

	if !exists {
		cmd := exec.Command("docker", "run",
			"--name", "sqlc_sqltest_docker_clickhouse",
			"-e", "CLICKHOUSE_DB=default",
			"-e", "CLICKHOUSE_USER=default",
			"-e", "CLICKHOUSE_PASSWORD=mysecretpassword",
			"-p", "9000:9000",
			"-p", "8123:8123",
			"-d",
			"clickhouse/clickhouse-server:latest",
		)

		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		msg := `Conflict. The container name "/sqlc_sqltest_docker_clickhouse" is already in use by container`
		if !strings.Contains(string(output), msg) && err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	// Create a ticker that fires every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout reached: %w", ctx.Err())

		case <-ticker.C:
			conn, err := sql.Open("clickhouse", uri)
			if err != nil {
				slog.Debug("sqltest", "connect", err)
				continue
			}
			defer conn.Close()
			if err := conn.PingContext(ctx); err != nil {
				slog.Debug("sqltest", "ping", err)
				continue
			}
			return uri, nil
		}
	}
}
