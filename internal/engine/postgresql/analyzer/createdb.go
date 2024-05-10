package analyzer

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (a *Analyzer) createDb(ctx context.Context, migrations []string) (string, error) {
	hash := a.fnv(migrations)
	name := fmt.Sprintf("sqlc_%s", hash)

	serverUri := a.replacer.Replace(a.servers[0].URI)
	pool, err := a.serverCache.Open(ctx, serverUri)
	if err != nil {
		return "", err
	}

	uri, err := url.Parse(serverUri)
	if err != nil {
		return "", err
	}
	uri.Path = name

	key := uri.String()
	_, err, _ = a.flight.Do(key, func() (interface{}, error) {
		// TODO: Use a parameterized query
		row := pool.QueryRow(ctx,
			fmt.Sprintf(`SELECT datname FROM pg_database WHERE datname = '%s'`, name))

		var datname string
		if err := row.Scan(&datname); err == nil {
			slog.Info("database exists", "name", name)
			return nil, nil
		}

		slog.Info("creating database", "name", name)
		if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, name)); err != nil {
			return nil, err
		}

		conn, err := pgx.Connect(ctx, uri.String())
		if err != nil {
			return nil, fmt.Errorf("connect %s: %s", name, err)
		}
		defer conn.Close(ctx)

		for _, q := range migrations {
			if len(strings.TrimSpace(q)) == 0 {
				continue
			}
			if _, err := conn.Exec(ctx, q); err != nil {
				return nil, fmt.Errorf("%s: %s", q, err)
			}
		}
		return nil, nil
	})

	if err != nil {
		return "", err
	}

	return key, err
}
