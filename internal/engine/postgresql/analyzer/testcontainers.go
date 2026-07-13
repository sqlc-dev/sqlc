package analyzer

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func startPostgresContainer(ctx context.Context, image string, migrations []string) (string, func(context.Context), error) {
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return "", nil, fmt.Errorf("testcontainers: disable ryuk: %w", err)
	}

	tmpDir, initScriptPaths, err := writeMigrationsTempDir(migrations)
	if err != nil {
		return "", nil, err
	}
	cleanupTmp := func() { os.RemoveAll(tmpDir) }

	opts := []testcontainers.ContainerCustomizer{
		tcpostgres.WithDatabase("sqlc"),
		tcpostgres.WithUsername("sqlc"),
		tcpostgres.WithPassword("sqlc"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30 * time.Second),
		),
	}
	if len(initScriptPaths) > 0 {
		opts = append(opts, tcpostgres.WithInitScripts(initScriptPaths...))
	}

	container, err := tcpostgres.Run(ctx, image, opts...)
	if err != nil {
		cleanupTmp()
		return "", nil, err
	}
	uri, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		container.Terminate(context.Background()) //nolint:errcheck
		cleanupTmp()
		return "", nil, err
	}

	return uri, func(ctx context.Context) {
		container.Terminate(ctx) //nolint:errcheck
		cleanupTmp()
	}, nil
}

func writeMigrationsTempDir(migrations []string) (dir string, paths []string, err error) {
	dir, err = os.MkdirTemp("", "sqlc-migrations-*")
	if err != nil {
		return "", nil, fmt.Errorf("testcontainers: create temp dir: %w", err)
	}
	for i, sql := range migrations {
		path := fmt.Sprintf("%s/%03d.sql", dir, i)
		if err := os.WriteFile(path, []byte(sql), 0600); err != nil {
			os.RemoveAll(dir)
			return "", nil, fmt.Errorf("testcontainers: write migration: %w", err)
		}
		paths = append(paths, path)
	}
	return dir, paths, nil
}
