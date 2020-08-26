// +build examples

package authors

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestAuthors(t *testing.T) {
	sdb, cleanup := sqltest.PostgreSQL(t, []string{"schema.sql"})
	defer cleanup()

	ctx := context.Background()
	db := New(sdb)

	// list all authors
	authors, err := db.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authors)

	// create an author
	insertedAuthor, err := db.CreateAuthor(ctx, CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(insertedAuthor)

	// get the author we just inserted
	fetchedAuthor, err := db.GetAuthor(ctx, insertedAuthor.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fetchedAuthor)
}
