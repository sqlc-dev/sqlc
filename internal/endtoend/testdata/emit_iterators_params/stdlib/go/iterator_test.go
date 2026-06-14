package querytest_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/sqlc-dev/sqlc/endtoend/emit_iterators_params/stdlib/go"
)

func TestIterAuthorsByMinIDFiltersRows(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE author (
			id INTEGER NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			bio TEXT
		);
		INSERT INTO author (id, name, bio) VALUES
			(1, 'ada', 'a'),
			(2, 'grace', 'b'),
			(3, 'linus', 'c');
	`); err != nil {
		t.Fatal(err)
	}

	q := querytest.New(db)

	slice, err := q.ListAuthorsByMinID(ctx, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(slice) != 2 {
		t.Fatalf("slice len=%d", len(slice))
	}

	var streamed []querytest.Author
	for author, err := range q.IterAuthorsByMinID(ctx, 2) {
		if err != nil {
			t.Fatal(err)
		}
		streamed = append(streamed, author)
	}
	if len(streamed) != len(slice) {
		t.Fatalf("iter len=%d slice len=%d", len(streamed), len(slice))
	}
	for i := range slice {
		if streamed[i] != slice[i] {
			t.Fatalf("row %d mismatch", i)
		}
	}
}
