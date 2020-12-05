// +build examples

package pggroup_test

import (
	"context"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/examples/group/postgresql/pggroup"
	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestBooks(t *testing.T) {
	db, cleanup := sqltest.PostgreSQL(t, []string{"../sql/schema/schema.sql"})
	defer cleanup()

	ctx := context.Background()

	dbauthors := pggroup.NewAuthors(db)

	// create an author
	a, err := dbauthors.Create(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}

	// create transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	txbooks := pggroup.NewBooks(tx)

	// save first book
	now := time.Now()
	_, err = txbooks.Create(ctx, pggroup.BooksCreateParams{
		AuthorID:  a.AuthorID,
		Isbn:      "1",
		Title:     "my book title",
		BookType:  pggroup.BookTypeFICTION,
		Year:      2016,
		Available: now,
		Tags:      []string{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// save second book
	b1, err := txbooks.Create(ctx, pggroup.BooksCreateParams{
		AuthorID:  a.AuthorID,
		Isbn:      "2",
		Title:     "the second book",
		BookType:  pggroup.BookTypeFICTION,
		Year:      2016,
		Available: now,
		Tags:      []string{"cool", "unique"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// update the title and tags
	err = txbooks.Update(ctx, pggroup.BooksUpdateParams{
		BookID: b1.BookID,
		Title:  "changed second title",
		Tags:   []string{"cool", "disastor"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// save third book
	_, err = txbooks.Create(ctx, pggroup.BooksCreateParams{
		AuthorID:  a.AuthorID,
		Isbn:      "3",
		Title:     "the third book",
		BookType:  pggroup.BookTypeFICTION,
		Year:      2001,
		Available: now,
		Tags:      []string{"cool"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// save fourth book
	b3, err := txbooks.Create(ctx, pggroup.BooksCreateParams{
		AuthorID:  a.AuthorID,
		Isbn:      "4",
		Title:     "4th place finisher",
		BookType:  pggroup.BookTypeNONFICTION,
		Year:      2011,
		Available: now,
		Tags:      []string{"other"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// tx commit
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	dbbooks := pggroup.NewBooks(db)

	// upsert, changing ISBN and title
	err = dbbooks.UpdateISBN(ctx, pggroup.BooksUpdateISBNParams{
		BookID: b3.BookID,
		Isbn:   "NEW ISBN",
		Title:  "never ever gonna finish, a quatrain",
		Tags:   []string{"someother"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	books0, err := dbbooks.ListByTitleYear(ctx, pggroup.BooksListByTitleYearParams{
		Title: "my book title",
		Year:  2016,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Format(time.RFC822Z))
		author, err := dbauthors.Get(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s\n", book.BookID, author.Name)
	}

	// find a book with either "cool" or "other" tag
	t.Logf("---------\nTag search results:\n")
	res, err := dbbooks.ListByTags(ctx, []string{"cool", "other", "someother"})
	if err != nil {
		t.Fatal(err)
	}
	for _, ab := range res {
		t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name, ab.Isbn, ab.Tags)
	}

	// get book 4 and delete
	b5, err := dbbooks.Get(ctx, b3.BookID)
	if err != nil {
		t.Fatal(err)
	}
	if err := dbbooks.Delete(ctx, b5.BookID); err != nil {
		t.Fatal(err)
	}
}
