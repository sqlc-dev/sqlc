package docker

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlSync sync.Once
var mysqlHost string

func StartMySQLServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}

	{
		_, err := exec.Command("docker", "pull", "mysql:8").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: mysql:8 %w", err)
		}
	}

	var syncErr error
	mysqlSync.Do(func() {
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
			syncErr = err
			return
		}

		// Create a ticker that fires every 10ms
		ticker := time.NewTicker(10 * time.Millisecond)
		defer ticker.Stop()

		uri := "root:mysecretpassword@/dinotest"

		db, err := sql.Open("mysql", uri)
		if err != nil {
			syncErr = fmt.Errorf("sql.Open: %w", err)
			return
		}

		for {
			select {
			case <-ctx.Done():
				syncErr = fmt.Errorf("timeout reached: %w", ctx.Err())
				return

			case <-ticker.C:
				// Run your function here
				if err := db.PingContext(ctx); err != nil {
					continue
				}
				mysqlHost = uri
				return
			}
		}
	})

	if syncErr != nil {
		return "", syncErr
	}

	if mysqlHost == "" {
		return "", fmt.Errorf("mysql server setup failed")
	}

	return mysqlHost, nil
}
