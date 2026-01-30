
package authors

import (
	"context"
	"testing"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
	_ "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/query"
)

func ptr(s string) *string {
	return &s
}

func TestAuthors(t *testing.T) {
	ctx := context.Background()
	db := local.YDB(t, []string{"schema.sql"})
	defer db.Close(ctx)

	q := New(db.Query())

	// list all authors
	authors, err := q.ListAuthors(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(authors)

	// create an author
	insertedAuthor, err := q.CreateOrUpdateAuthor(ctx, CreateOrUpdateAuthorParams{
		Name: "Brian Kernighan",
		Bio:  ptr("Co-author of The C Programming Language and The Go Programming Language"),
	}, query.WithIdempotent())
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

	// drop table
	err = q.DropTable(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
