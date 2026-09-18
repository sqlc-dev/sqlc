//go:build examples

package authors

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/duckdb/duckdb-go/v2"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()

	// DuckDB runs in the process: an empty data source is a database in
	// memory that lives as long as the connection.
	sdb, err := sql.Open("duckdb", "")
	if err != nil {
		t.Fatal(err)
	}
	defer sdb.Close()

	schema, err := os.ReadFile("schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := sdb.ExecContext(ctx, string(schema)); err != nil {
		t.Fatal(err)
	}

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
