package docker

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
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
		_, err := exec.Command("docker", "pull", "clickhouse:lts").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: clickhouse:lts %w", err)
		}
	}

	var exists bool
	{
		cmd := exec.Command("docker", "container", "inspect", "sqlc_sqltest_docker_clickhouse")
		// This means we've already started the container
		exists = cmd.Run() == nil
	}

	if !exists {
		cmd := exec.Command("docker", "run",
			"--name", "sqlc_sqltest_docker_clickhouse",
			"-p", "9000:9000",
			"-p", "8123:8123",
			"-e", "CLICKHOUSE_DB=default",
			"-e", "CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1",
			"-d",
			"clickhouse:lts",
		)

		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		msg := `Conflict. The container name "/sqlc_sqltest_docker_clickhouse" is already in use by container`
		if !strings.Contains(string(output), msg) && err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Create a ticker that fires every 10ms
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	// ClickHouse DSN format: clickhouse://host:port
	dsn := "clickhouse://localhost:9000/default"

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout reached: %w", ctx.Err())

		case <-ticker.C:
			db, err := sql.Open("clickhouse", dsn)
			if err != nil {
				slog.Debug("sqltest", "open", err)
				continue
			}

			if err := db.PingContext(ctx); err != nil {
				slog.Debug("sqltest", "ping", err)
				db.Close()
				continue
			}

			db.Close()
			return dsn, nil
		}
	}
}
