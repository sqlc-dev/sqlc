// +build examples

package jsonemptyarrays_test

import (
	"context"
	"encoding/json"
	"github.com/kyleconroy/sqlc/examples/booktest/postgresql/jsonemptyarrays"
	"github.com/kyleconroy/sqlc/internal/sqltest"
	"testing"
)

func TestBooks(t *testing.T) {
	db, cleanup := sqltest.PostgreSQL(t, []string{"../schema.sql"})
	defer cleanup()

	ctx := context.Background()
	dq := jsonemptyarrays.New(db)

	// lookup books with no results
	books, err := dq.BooksByTitleYear(ctx, jsonemptyarrays.BooksByTitleYearParams{
		Title: "my book title",
		Year:  2016,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 0 {
		t.Fatal("books should be empty")
	}

	// assert json encoding returns empty array
	data, err := json.Marshal(&books)
	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "[]" {
		t.Fatalf("json.Marshal should encode an empty array got: %s", string(data))
	}
}
