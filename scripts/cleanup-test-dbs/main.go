package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

const query = `
SELECT datname
FROM pg_database
WHERE datname LIKE 'sqlc_test_%'
`

func run() error {
	ctx := context.Background()
	dburi := os.Getenv("POSTGRESQL_SERVER_URI")
	if dburi == "" {
		return fmt.Errorf("POSTGRESQL_SERVER_URI is empty")
	}
	pool, err := pgxpool.New(ctx, dburi)
	if err != nil {
		return err
	}

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return err
	}

	names, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return err
	}

	for _, name := range names {
		drop := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name)
		if _, err := pool.Exec(ctx, drop); err != nil {
			return err
		}
		log.Println("dropping database", name)
	}

	return nil
}
