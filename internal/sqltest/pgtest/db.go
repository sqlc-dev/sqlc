package pgtest

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func id() string {
	b := make([]rune, 10)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CreateDatabase(t *testing.T, ctx context.Context, pool *pgxpool.Pool) (string, func()) {
	t.Helper()

	name := fmt.Sprintf("sqlc_test_%s", id())

	if _, err := pool.Exec(ctx, fmt.Sprintf(`CREATE DATABASE "%s"`, name)); err != nil {
		t.Fatal(err)
	}

	dropQuery := fmt.Sprintf(`DROP DATABASE IF EXISTS "%s" WITH (FORCE)`, name)

	return name, func() {
		if _, err := pool.Exec(ctx, dropQuery); err != nil {
			t.Fatal(err)
		}
	}
}
