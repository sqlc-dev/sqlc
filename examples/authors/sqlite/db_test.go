//go:build examples

package authors

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest"
)

func TestAuthors(t *testing.T) {
	sdb, cleanup := sqltest.SQLite(t, []string{"schema.sql"})
	defer sdb.Close()
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
	result, err := db.CreateAuthor(ctx, CreateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
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
	t.Log(fetchedAuthor)
}

func TestIterAuthors(t *testing.T) {
	sdb, cleanup := sqltest.SQLite(t, []string{"schema.sql"})
	defer sdb.Close()
	defer cleanup()

	ctx := context.Background()
	db := New(sdb)

	for _, name := range []string{"Rob Pike", "Brian Kernighan", "Ken Thompson"} {
		if _, err := db.CreateAuthor(ctx, CreateAuthorParams{Name: name}); err != nil {
			t.Fatal(err)
		}
	}

	// the query only runs once the iterator is ranged over
	authors, ferr := db.IterAuthors(ctx)
	if inUse := sdb.Stats().InUse; inUse != 0 {
		t.Fatalf("%d connections in use before ranging, want 0", inUse)
	}
	var names []string
	for author := range authors {
		names = append(names, author.Name)
	}
	if err := ferr(); err != nil {
		t.Fatal(err)
	}
	if want := []string{"Brian Kernighan", "Ken Thompson", "Rob Pike"}; !slices.Equal(names, want) {
		t.Fatalf("got %v, want %v", names, want)
	}

	// the iterator is single-use: ranging over it again yields nothing and
	// reports an error
	for author := range authors {
		t.Errorf("second range yielded %v", author)
	}
	if err := ferr(); err == nil {
		t.Error("ranging twice reported no error")
	}

	// breaking out of the loop early releases the connection
	authors, ferr = db.IterAuthors(ctx)
	for range authors {
		break
	}
	if err := ferr(); err != nil {
		t.Fatal(err)
	}
	if inUse := sdb.Stats().InUse; inUse != 0 {
		t.Errorf("%d connections in use after breaking out, want 0", inUse)
	}

	// a query that fails reports its error through the returned func
	cctx, cancel := context.WithCancel(ctx)
	cancel()
	authors, ferr = db.IterAuthors(cctx)
	for author := range authors {
		t.Errorf("cancelled query yielded %v", author)
	}
	if err := ferr(); !errors.Is(err, context.Canceled) {
		t.Errorf("got error %v, want context.Canceled", err)
	}
}
