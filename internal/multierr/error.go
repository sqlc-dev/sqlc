package multierr

import (
	"fmt"

	"github.com/kyleconroy/sqlc/internal/pg"
	"github.com/kyleconroy/sqlc/internal/source"
)

type FileError struct {
	Filename string
	Line     int
	Column   int
	Err      error
}

func (e *FileError) Unwrap() error {
	return e.Err
}

type Error struct {
	errs []*FileError
}

func (e *Error) Add(filename, in string, loc int, err error) {
	line := 1
	column := 1
	if lerr, ok := err.(pg.Error); ok {
		if lerr.Location != 0 {
			loc = lerr.Location
		}
	}
	if in != "" && loc != 0 {
		line, column = source.LineNumber(in, loc)
	}
	e.errs = append(e.errs, &FileError{filename, line, column, err})
}

func (e *Error) Errs() []*FileError {
	return e.errs
}

func (e *Error) Error() string {
	return fmt.Sprintf("multiple errors: %d errors", len(e.errs))
}

func New() *Error {
	return &Error{}
}
