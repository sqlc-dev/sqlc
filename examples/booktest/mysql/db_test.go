// +build examples

package booktest

import (
	"context"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestBooks(t *testing.T) {
	db, cleanup := sqltest.MySQL(t, []string{"schema.sql"})
	defer cleanup()

	ctx := context.Background()
	dq := New(db)

	// create an author
	result, err := dq.CreateAuthor(ctx, "Unknown Master")
	if err != nil {
		t.Fatal(err)
	}
	authorID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// create transaction
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}

	tq := dq.WithTx(tx)

	// save first book
	now := time.Now()
	_, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  int32(authorID),
		Isbn:      "1",
		Title:     "my book title",
		BookType:  BookTypeFICTION,
		Yr:        2016,
		Available: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	// save second book
	result, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  int32(authorID),
		Isbn:      "2",
		Title:     "the second book",
		BookType:  BookTypeFICTION,
		Yr:        2016,
		Available: now,
		Tags:      "cool,unique",
	})
	if err != nil {
		t.Fatal(err)
	}
	bookOneID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// update the title and tags
	err = tq.UpdateBook(ctx, UpdateBookParams{
		BookID: int32(bookOneID),
		Title:  "changed second title",
		Tags:   "cool,disastor",
	})
	if err != nil {
		t.Fatal(err)
	}

	// save third book
	_, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  int32(authorID),
		Isbn:      "3",
		Title:     "the third book",
		BookType:  BookTypeFICTION,
		Yr:        2001,
		Available: now,
		Tags:      "cool",
	})
	if err != nil {
		t.Fatal(err)
	}

	// save fourth book
	result, err = tq.CreateBook(ctx, CreateBookParams{
		AuthorID:  int32(authorID),
		Isbn:      "4",
		Title:     "4th place finisher",
		BookType:  BookTypeNONFICTION,
		Yr:        2011,
		Available: now,
		Tags:      "other",
	})
	if err != nil {
		t.Fatal(err)
	}
	bookThreeID, err := result.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}

	// tx commit
	err = tx.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// upsert, changing ISBN and title
	err = dq.UpdateBookISBN(ctx, UpdateBookISBNParams{
		BookID: int32(bookThreeID),
		Isbn:   "NEW ISBN",
		Title:  "never ever gonna finish, a quatrain",
		Tags:   "someother",
	})
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	books0, err := dq.BooksByTitleYear(ctx, BooksByTitleYearParams{
		Title: "my book title",
		Yr:    2016,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, book := range books0 {
		t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Format(time.RFC822Z))
		author, err := dq.GetAuthor(ctx, book.AuthorID)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Book %d author: %s\n", book.BookID, author.Name)
	}

	// find a book with either "cool" or "other" tag
	t.Logf("---------\nTag search results:\n")
	res, err := dq.BooksByTags(ctx, "cool")
	if err != nil {
		t.Fatal(err)
	}
	for _, ab := range res {
		t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name, ab.Isbn, ab.Tags)
	}

	// TODO: call say_hello(varchar)

	// get book 4 and delete
	b5, err := dq.GetBook(ctx, int32(bookThreeID))
	if err != nil {
		t.Fatal(err)
	}
	if err := dq.DeleteBook(ctx, b5.BookID); err != nil {
		t.Fatal(err)
	}

}
