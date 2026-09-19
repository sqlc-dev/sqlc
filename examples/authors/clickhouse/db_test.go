//go:build examples

package authors

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/sqlc-dev/sqlc/examples/internal/local"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()
	uri := local.ClickHouse(t, []string{"schema.sql"})
	sdb, err := sql.Open("clickhouse", uri)
	if err != nil {
		t.Fatal(err)
	}
	defer sdb.Close()

	db := New(sdb)

	// list all authors
	authors, err := db.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authors)

	// create an author
	err = db.CreateAuthor(ctx, CreateAuthorParams{
		ID:   1,
		Name: "Brian Kernighan",
		Bio:  sql.NullString{String: "Co-author of The C Programming Language and The Go Programming Language", Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	// get the author we just inserted
	fetchedAuthor, err := db.GetAuthor(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if fetchedAuthor.Name != "Brian Kernighan" || !fetchedAuthor.Bio.Valid {
		t.Fatalf("unexpected author: %+v", fetchedAuthor)
	}
	t.Log(fetchedAuthor)

	// delete the author
	if err := db.DeleteAuthor(ctx, 1); err != nil {
		t.Fatal(err)
	}
}
