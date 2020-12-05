// +build examples

package mygroupiface_test

import (
	"context"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/examples/group/mysql/mygroupiface"
	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestBooks(t *testing.T) {
	db, cleanup := sqltest.MySQL(t, []string{"../sql/schema/schema.sql"})
	defer cleanup()

	ctx := context.Background()

	dbauthors := mygroupiface.NewAuthors(db)

	// create an author
	result, err := dbauthors.Create(ctx, "Unknown Master")
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

	txbooks := mygroupiface.NewBooks(tx)

	// save first book
	now := time.Now()
	_, err = txbooks.Create(ctx, mygroupiface.BooksCreateParams{
		AuthorID:  int32(authorID),
		Isbn:      "1",
		Title:     "my book title",
		BookType:  mygroupiface.BooksBookTypeFICTION,
		Yr:        2016,
		Available: now,
	})
	if err != nil {
		t.Fatal(err)
	}

	// save second book
	result, err = txbooks.Create(ctx, mygroupiface.BooksCreateParams{
		AuthorID:  int32(authorID),
		Isbn:      "2",
		Title:     "the second book",
		BookType:  mygroupiface.BooksBookTypeFICTION,
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
	err = txbooks.Update(ctx, mygroupiface.BooksUpdateParams{
		BookID: int32(bookOneID),
		Title:  "changed second title",
		Tags:   "cool,disastor",
	})
	if err != nil {
		t.Fatal(err)
	}

	// save third book
	_, err = txbooks.Create(ctx, mygroupiface.BooksCreateParams{
		AuthorID:  int32(authorID),
		Isbn:      "3",
		Title:     "the third book",
		BookType:  mygroupiface.BooksBookTypeFICTION,
		Yr:        2001,
		Available: now,
		Tags:      "cool",
	})
	if err != nil {
		t.Fatal(err)
	}

	// save fourth book
	result, err = txbooks.Create(ctx, mygroupiface.BooksCreateParams{
		AuthorID:  int32(authorID),
		Isbn:      "4",
		Title:     "4th place finisher",
		BookType:  mygroupiface.BooksBookTypeNONFICTION,
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

	dbbooks := mygroupiface.NewBooks(db)

	// upsert, changing ISBN and title
	err = dbbooks.UpdateISBN(ctx, mygroupiface.BooksUpdateISBNParams{
		BookID: int32(bookThreeID),
		Isbn:   "NEW ISBN",
		Title:  "never ever gonna finish, a quatrain",
		Tags:   "someother",
	})
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	books0, err := dbbooks.ListByTitleYear(ctx, mygroupiface.BooksListByTitleYearParams{
		Title: "my book title",
		Yr:    2016,
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
	res, err := dbbooks.ListByTags(ctx, "cool")
	if err != nil {
		t.Fatal(err)
	}
	for _, ab := range res {
		t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name, ab.Isbn, ab.Tags)
	}

	// get book 4 and delete
	b5, err := dbbooks.Get(ctx, int32(bookThreeID))
	if err != nil {
		t.Fatal(err)
	}
	if err := dbbooks.Delete(ctx, b5.BookID); err != nil {
		t.Fatal(err)
	}
}
