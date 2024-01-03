//go:build examples
// +build examples

package authors

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func TestAuthors(t *testing.T) {
	ctx := context.Background()
	uri := local.MySQL(t, []string{"schema.sql"})
	sdb, err := sql.Open("mysql", uri)
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
