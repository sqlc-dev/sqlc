//go:build examples
// +build examples

package authors

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"

	"github.com/sqlc-dev/sqlc/internal/sqltest/hosted"
)

func TestAuthors(t *testing.T) {
	uri := hosted.PostgreSQL(t, []string{"schema.sql"})
	db, err := sql.Open("postgres", uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	q := New(db)

	// list all authors
	authors, err := q.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authors)

	// create an author
	insertedAuthor, err := q.CreateAuthor(ctx, CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := q.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fetchedAuthor)
}
