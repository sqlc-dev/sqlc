// +build examples

package authors

import (
	"context"
	"database/sql"
	"testing"

	"github.com/kyleconroy/sqlc/internal/sqltest"
)

// An example of an application-level structure
type AppAuthor struct {
	ID   int64
	Name string
	Bio  string
}

// A mapper that uses a face interface
func MapToAuthor(a interface {
	GetID() int64
	GetName() string
	GetBio() sql.NullString
}) AppAuthor {
	return AppAuthor{
		ID:   a.GetID(),
		Name: a.GetName(),
		Bio:  a.GetBio().String,
	}
}

func TestAuthors(t *testing.T) {
	sdb, cleanup := sqltest.PostgreSQL(t, "schema.sql")
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
	t.Log(MapToAuthor(&fetchedAuthor))
}
