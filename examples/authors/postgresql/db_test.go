//go:build examples
// +build examples

package authors

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()
	uri := local.PostgreSQL(t, []string{"schema.sql"})
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)

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
		Bio:  pgtype.Text{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
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
