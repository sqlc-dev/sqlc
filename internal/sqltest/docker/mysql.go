package docker

import (
	"context"
	"database/sql"
	"fmt"
	"os/exec"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlHost string

func StartMySQLServer(c context.Context) (string, error) {
	if err := Installed(); err != nil {
		return "", err
	}
	if mysqlHost != "" {
		return mysqlHost, nil
	}
	value, err, _ := flight.Do("mysql", func() (interface{}, error) {
		host, err := startMySQLServer(c)
		if err != nil {
			return "", err
		}
		mysqlHost = host
		return host, nil
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

func startMySQLServer(c context.Context) (string, error) {
	{
		_, err := exec.Command("docker", "pull", "mysql:9").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("docker pull: mysql:9 %w", err)
		}
	}

	var exists bool
	{
		cmd := exec.Command("docker", "container", "inspect", "sqlc_sqltest_docker_mysql")
		// This means we've already started the container
		exists = cmd.Run() == nil
	}

	if !exists {
		cmd := exec.Command("docker", "run",
			"--name", "sqlc_sqltest_docker_mysql",
			"-e", "MYSQL_ROOT_PASSWORD=mysecretpassword",
			"-e", "MYSQL_DATABASE=dinotest",
			"-p", "3306:3306",
			"-d",
			"mysql:9",
		)

		output, err := cmd.CombinedOutput()
		fmt.Println(string(output))

		msg := `Conflict. The container name "/sqlc_sqltest_docker_mysql" is already in use by container`
		if !strings.Contains(string(output), msg) && err != nil {
			return "", err
		}
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	// Create a ticker that fires every 10ms
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	uri := "root:mysecretpassword@/dinotest?multiStatements=true&parseTime=true"

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return "", fmt.Errorf("sql.Open: %w", err)
	}

	defer db.Close()

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
