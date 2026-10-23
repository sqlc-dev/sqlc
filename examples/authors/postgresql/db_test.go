//go:build examples

package authors

import (
	"context"
	"errors"
	"slices"
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

func TestIterAuthors(t *testing.T) {
	ctx := context.Background()
	uri := local.PostgreSQL(t, []string{"schema.sql"})
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)

	q := New(db)

	var last Author
	for _, name := range []string{"Rob Pike", "Brian Kernighan", "Ken Thompson"} {
		if last, err = q.CreateAuthor(ctx, CreateAuthorParams{Name: name}); err != nil {
			t.Fatal(err)
		}
	}

	// the query only runs once the iterator is ranged over, so the
	// connection is still free for other queries until then
	authors, ferr := q.IterAuthors(ctx)
	if _, err := q.GetAuthor(ctx, last.ID); err != nil {
		t.Fatalf("query before ranging: %v", err)
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
	authors, ferr = q.IterAuthors(ctx)
	for range authors {
		break
	}
	if err := ferr(); err != nil {
		t.Fatal(err)
	}
	if _, err := q.GetAuthor(ctx, last.ID); err != nil {
		t.Errorf("query after breaking out: %v", err)
	}

	// a query that fails reports its error through the returned func
	cctx, cancel := context.WithCancel(ctx)
	cancel()
	authors, ferr = q.IterAuthors(cctx)
	for author := range authors {
		t.Errorf("cancelled query yielded %v", author)
	}
	if err := ferr(); !errors.Is(err, context.Canceled) {
		t.Errorf("got error %v, want context.Canceled", err)
	}
}
