// +build examples

package pggroupstmtiface_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/kyleconroy/sqlc/examples/group/postgresql/pggroupstmtiface"
	"github.com/kyleconroy/sqlc/internal/sqltest"
)

func TestBooks(t *testing.T) {
	db, cleanup := sqltest.PostgreSQL(t, []string{"../sql/schema/schema.sql"})
	defer cleanup()

	ctx := context.Background()

	// create an author
	a, err := func() (retModel pggroupstmtiface.Author, retErr error) {
		dbauthors := pggroupstmtiface.NewAuthors(db)
		dbcreate, err := dbauthors.PrepareCreate(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbcreate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		a, err := dbcreate.Exec(ctx, "Unknown Master")
		if err != nil {
			retErr = err
			return
		}
		retModel = a
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// create transaction
	b3, err := func() (retModel pggroupstmtiface.Book, retErr error) {
		tx, err := db.Begin()
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := tx.Rollback(); err != nil && err != sql.ErrTxDone && retErr == nil {
				retErr = err
			}
		}()

		txbooks := pggroupstmtiface.NewBooks(tx)

		txcreate, err := txbooks.PrepareCreate(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := txcreate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()

		// save first book
		now := time.Now()
		_, err = txcreate.Exec(ctx, pggroupstmtiface.BooksCreateParams{
			AuthorID:  a.AuthorID,
			Isbn:      "1",
			Title:     "my book title",
			BookType:  pggroupstmtiface.BookTypeFICTION,
			Year:      2016,
			Available: now,
			Tags:      []string{},
		})
		if err != nil {
			retErr = err
			return
		}

		// save second book
		b1, err := txcreate.Exec(ctx, pggroupstmtiface.BooksCreateParams{
			AuthorID:  a.AuthorID,
			Isbn:      "2",
			Title:     "the second book",
			BookType:  pggroupstmtiface.BookTypeFICTION,
			Year:      2016,
			Available: now,
			Tags:      []string{"cool", "unique"},
		})
		if err != nil {
			retErr = err
			return
		}

		// update the title and tags
		txupdate, err := txbooks.PrepareUpdate(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := txupdate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		err = txupdate.Exec(ctx, pggroupstmtiface.BooksUpdateParams{
			BookID: b1.BookID,
			Title:  "changed second title",
			Tags:   []string{"cool", "disastor"},
		})
		if err != nil {
			retErr = err
			return
		}

		// save third book
		_, err = txcreate.Exec(ctx, pggroupstmtiface.BooksCreateParams{
			AuthorID:  a.AuthorID,
			Isbn:      "3",
			Title:     "the third book",
			BookType:  pggroupstmtiface.BookTypeFICTION,
			Year:      2001,
			Available: now,
			Tags:      []string{"cool"},
		})
		if err != nil {
			retErr = err
			return
		}

		// save fourth book
		b3, err := txcreate.Exec(ctx, pggroupstmtiface.BooksCreateParams{
			AuthorID:  a.AuthorID,
			Isbn:      "4",
			Title:     "4th place finisher",
			BookType:  pggroupstmtiface.BookTypeNONFICTION,
			Year:      2011,
			Available: now,
			Tags:      []string{"other"},
		})
		if err != nil {
			retErr = err
			return
		}

		// tx commit
		err = tx.Commit()
		if err != nil {
			retErr = err
			return
		}
		retModel = b3
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// upsert, changing ISBN and title
	err = func() (retErr error) {
		dbbooks := pggroupstmtiface.NewBooks(db)
		dbupdate, err := dbbooks.PrepareUpdateISBN(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbupdate.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		retErr = dbupdate.Exec(ctx, pggroupstmtiface.BooksUpdateISBNParams{
			BookID: b3.BookID,
			Isbn:   "NEW ISBN",
			Title:  "never ever gonna finish, a quatrain",
			Tags:   []string{"someother"},
		})
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// retrieve first book
	err = func() (retErr error) {
		dbbooks := pggroupstmtiface.NewBooks(db)
		dblist, err := dbbooks.PrepareListByTitleYear(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dblist.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		books0, err := dblist.Exec(ctx, pggroupstmtiface.BooksListByTitleYearParams{
			Title: "my book title",
			Year:  2016,
		})
		if err != nil {
			retErr = err
			return
		}
		dbauthors := pggroupstmtiface.NewAuthors(db)
		dbget, err := dbauthors.PrepareGet(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbget.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		for _, book := range books0 {
			t.Logf("Book %d (%s): %s available: %s\n", book.BookID, book.BookType, book.Title, book.Available.Format(time.RFC822Z))
			author, err := dbget.Exec(ctx, book.AuthorID)
			if err != nil {
				retErr = err
				return
			}
			t.Logf("Book %d author: %s\n", book.BookID, author.Name)
		}
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// find a book with either "cool" or "other" tag
	err = func() (retErr error) {
		t.Logf("---------\nTag search results:\n")
		dbbooks := pggroupstmtiface.NewBooks(db)
		dblist, err := dbbooks.PrepareListByTags(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dblist.Close(); err != nil && retErr == nil {
				retErr = err
				return
			}
		}()
		res, err := dblist.Exec(ctx, []string{"cool", "other", "someother"})
		if err != nil {
			retErr = err
			return
		}
		for _, ab := range res {
			t.Logf("Book %d: '%s', Author: '%s', ISBN: '%s' Tags: '%v'\n", ab.BookID, ab.Title, ab.Name, ab.Isbn, ab.Tags)
		}
		return
	}()
	if err != nil {
		t.Fatal(err)
	}

	// get book 4 and delete
	err = func() (retErr error) {
		dbbooks := pggroupstmtiface.NewBooks(db)
		dbget, err := dbbooks.PrepareGet(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbget.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		b5, err := dbget.Exec(ctx, b3.BookID)
		if err != nil {
			retErr = err
			return
		}
		dbdelete, err := dbbooks.PrepareDelete(ctx)
		if err != nil {
			retErr = err
			return
		}
		defer func() {
			if err := dbdelete.Close(); err != nil && retErr == nil {
				retErr = err
			}
		}()
		if err := dbdelete.Exec(ctx, b5.BookID); err != nil {
			retErr = err
			return
		}
		return
	}()
	if err != nil {
		t.Fatal(err)
	}
}
