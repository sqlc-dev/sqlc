//go:build examples
// +build examples

package update

import (
	"context"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestAuthor(t *testing.T) {
	sdb, cleanup := sqltest.MySQL(t, []string{"schema.sql"})
	defer cleanup()

	ctx := context.Background()
	db := New(sdb)

	// create an author
	result, err := db.CreateAuthor(ctx, CreateAuthorParams{
		Name:      "Brian Kernighan",
		DeletedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		t.Fatal(err)
	}
	authorID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authorID)

	// get the author we just inserted
	fetchedAuthor, err := db.GetAuthor(ctx, authorID)
	if err != nil {
		t.Fatal(err)
	}

	// create a book
	_, err = db.CreateBook(ctx, true)
	if err != nil {
		t.Fatal(err)
	}

	err = db.DeleteAuthor(ctx, "Brian Kernighan")
	if err != nil {
		t.Fatal(err)
	}

	// get the author we just inserted
	newFetchedAuthor, err := db.GetAuthor(ctx, authorID)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fetchedAuthor)
	if newFetchedAuthor.DeletedAt.Unix() != fetchedAuthor.DeletedAt.Unix() && newFetchedAuthor.DeletedAt.Unix() != newFetchedAuthor.UpdatedAt.Unix() {
		t.Fatal("update fail")
	}
}
