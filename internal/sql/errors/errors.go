package errors

import (
	"errors"
	"fmt"
)

var AlreadyExists = errors.New("already exists")
var NotFound = errors.New("does not exist")

type Error struct {
	Err      error
	Code     string
	Message  string
	Location int
	// Hint     string
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s %s", e.Message, e.Err.Error())
}

func ColumnAlreadyExists(rel, col string) *Error {
	return &Error{
		Err:     AlreadyExists,
		Code:    "42701",
		Message: fmt.Sprintf("column \"%s\" of relation \"%s\"", col, rel),
	}
}

func ColumnNotFound(rel, col string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42703",
		Message: fmt.Sprintf("column \"%s\" of relation \"%s\"", col, rel),
	}
}

func RelationAlreadyExists(rel string) *Error {
	return &Error{
		Err:     AlreadyExists,
		Code:    "42P07",
		Message: fmt.Sprintf("relation \"%s\"", rel),
	}
}

func RelationNotFound(rel string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42P01",
		Message: fmt.Sprintf("relation \"%s\"", rel),
	}
}

func SchemaAlreadyExists(name string) *Error {
	return &Error{
		Err:     AlreadyExists,
		Code:    "42P06",
		Message: fmt.Sprintf("schema \"%s\"", name),
	}
}

func SchemaNotFound(sch string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "3F000",
		Message: fmt.Sprintf("schema \"%s\"", sch),
	}
}

func TypeAlreadyExists(typ string) *Error {
	return &Error{
		Err:     AlreadyExists,
		Code:    "42710",
		Message: fmt.Sprintf("type \"%s\"", typ),
	}
}

func TypeNotFound(typ string) *Error {
	return &Error{
		Err:     NotFound,
		Code:    "42704",
		Message: fmt.Sprintf("type \"%s\"", typ),
	}
}
