package docker

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func StartMySQLServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	cmd := exec.Command("docker", "run",
		"--name", "sqlc_sqltest_docker_mysql",
		"-e", "MYSQL_ROOT_PASSWORD=mysecretpassword",
		"-e", "MYSQL_DATABASE=dinotest",
		"-p", "3306:3306",
		"-d",
		"mysql:8",
	)

	output, err := cmd.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return "", err
	}

	// Create a ticker that fires every 10ms
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	uri := "root:mysecretpassword@/dinotest"

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return "", fmt.Errorf("sql.Open: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout reached: %w", ctx.Err())

		case <-ticker.C:
			// Run your function here
			if err := db.PingContext(ctx); err != nil {
				continue
			}
			return uri, nil
		}
	}
}
