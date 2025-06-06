// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package batch

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type BookType string

const (
	BookTypeFICTION    BookType = "FICTION"
	BookTypeNONFICTION BookType = "NONFICTION"
)

func (e *BookType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = BookType(s)
	case string:
		*e = BookType(s)
	default:
		return fmt.Errorf("unsupported scan type for BookType: %T", src)
	}
	return nil
}

type NullBookType struct {
	BookType BookType `json:"book_type"`
	Valid    bool     `json:"valid"` // Valid is true if BookType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullBookType) Scan(value interface{}) error {
	if value == nil {
		ns.BookType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.BookType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullBookType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.BookType), nil
}

type Author struct {
	AuthorID  int32  `json:"author_id"`
	Name      string `json:"name"`
	Biography []byte `json:"biography"`
}

type Book struct {
	BookID    int32              `json:"book_id"`
	AuthorID  int32              `json:"author_id"`
	Isbn      string             `json:"isbn"`
	BookType  BookType           `json:"book_type"`
	Title     string             `json:"title"`
	Year      int32              `json:"year"`
	Available pgtype.Timestamptz `json:"available"`
	Tags      []string           `json:"tags"`
}
