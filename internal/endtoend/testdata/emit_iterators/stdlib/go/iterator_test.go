package querytest_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"github.com/sqlc-dev/sqlc/endtoend/emit_iterators/stdlib/go"
)

func TestIterAuthorsMatchesListAuthors(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db := openTestDB(t, ctx)
	q := querytest.New(db)

	slice, err := q.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var streamed []querytest.Author
	for author, err := range q.IterAuthors(ctx) {
		if err != nil {
			t.Fatal(err)
		}
		streamed = append(streamed, author)
	}
	if len(streamed) != len(slice) {
		t.Fatalf("len(iter)=%d len(slice)=%d", len(streamed), len(slice))
	}
	for i := range slice {
		if streamed[i] != slice[i] {
			t.Fatalf("row %d: iter=%+v slice=%+v", i, streamed[i], slice[i])
		}
	}
}

func TestIterAuthorsLazyStart(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db := openTestDB(t, ctx)
	q := querytest.New(db)

	seq := q.IterAuthors(ctx)
	if seq == nil {
		t.Fatal("expected non-nil iterator")
	}

	closed, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	_ = closed.Close()
	lazyQ := querytest.New(closed)
	lazySeq := lazyQ.IterAuthors(ctx)
	gotErr := false
	for _, err := range lazySeq {
		if err != nil {
			gotErr = true
			break
		}
	}
	if !gotErr {
		t.Fatal("expected error when iterating with closed db")
	}
}

func TestIterAuthorsEarlyBreak(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db := openTestDB(t, ctx)
	q := querytest.New(db)

	count := 0
	for author, err := range q.IterAuthors(ctx) {
		if err != nil {
			t.Fatal(err)
		}
		if author.Name == "" {
			t.Fatal("empty name")
		}
		count++
		break
	}
	if count != 1 {
		t.Fatalf("count=%d", count)
	}

	// Connection must remain usable after early break.
	if _, err := q.GetAuthorByID(ctx, 1); err != nil {
		t.Fatalf("query after early break failed: %v", err)
	}
}

func openTestDB(t *testing.T, ctx context.Context) *sql.DB {
	t.Helper()
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
			(1, 'ada', 'first'),
			(2, 'grace', 'second');
	`); err != nil {
		t.Fatal(err)
	}
	return db
}
