//go:build examples

package authors

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/sqlc-dev/sqlc/examples/internal/local"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()
	uri := local.MSSQL(t, []string{"schema.sql"})
	sdb, err := sql.Open("sqlserver", uri)
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
	if fetchedAuthor.Name != "Brian Kernighan" || !fetchedAuthor.Bio.Valid {
		t.Fatalf("unexpected author: %+v", fetchedAuthor)
	}
	t.Log(fetchedAuthor)

	// delete the author
	if err := db.DeleteAuthor(ctx, insertedAuthor.ID); err != nil {
		t.Fatal(err)
	}
}
