//go:build examples
// +build examples

package booktest

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/sqlc-dev/sqlc/internal/sqltest/local"
)

func TestBooks(t *testing.T) {
	ctx := context.Background()
	uri := local.PostgreSQL(t, []string{"schema.sql"})
	db, err := pgx.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close(ctx)

	dq := New(db)

	// create an author
	a, err := dq.CreateAuthor(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}

	// create transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tq := dq.WithTx(tx)

	// save first book
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	_, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "1",
		Title:     "my book title",
		BookType:  BookTypeFICTION,
		Year:      2016,
		Available: now,
		Tags:      []string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// save second book
	b1, err := tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "2",
		Title:     "the second book",
		BookType:  BookTypeFICTION,
		Year:      2016,
		Available: now,
		Tags:      []string{"cool", "unique"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// update the title and tags
	err = tq.UpdateBook(ctx, UpdateBookParams{
		BookID: b1.BookID,
		Title:  "changed second title",
		Tags:   []string{"cool", "disastor"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// save third book
	_, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "3",
		Title:     "the third book",
		BookType:  BookTypeFICTION,
		Year:      2001,
		Available: now,
		Tags:      []string{"cool"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// save fourth book
	b3, err := tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  a.AuthorID,
		Isbn:      "4",
		Title:     "4th place finisher",
		BookType:  BookTypeNONFICTION,
		Year:      2011,
		Available: now,
		Tags:      []string{"other"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// tx commit
	err = tx.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// upsert, changing ISBN and title
	err = dq.UpdateBookISBN(ctx, UpdateBookISBNParams{
		BookID: b3.BookID,
		Isbn:   "NEW ISBN",
		Title:  "never ever gonna finish, a quatrain",
		Tags:   []string{"someother"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	books0, err := dq.BooksByTitleYear(ctx, BooksByTitleYearParams{
		Title: "my book title",
		Year:  2016,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Time.Format(time.RFC822Z))
		author, err := dq.GetAuthor(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s\n", book.BookID, author.Name)
	}

	// find a book with either "cool" or "other" tag
	t.Logf("---------\nTag search results:\n")
	res, err := dq.BooksByTags(ctx, []string{"cool", "other", "someother"})
	if err != nil {
		t.Fatal(err)
	}
	for _, ab := range res {
		t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name.String, ab.Isbn, ab.Tags)
	}

	// TODO: call say_hello(varchar)

	// get book 4 and delete
	b5, err := dq.GetBook(ctx, b3.BookID)
	if err != nil {
		t.Fatal(err)
	}
	if err := dq.DeleteBook(ctx, b5.BookID); err != nil {
		t.Fatal(err)
	}

}
