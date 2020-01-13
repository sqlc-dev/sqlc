package pg

import (
	"fmt"
)

type Error struct {
	Message  string
	Code     string
	Hint     string
	Location int
}

func (e Error) Error() string {
	return e.Message
}

func ErrorColumnAlreadyExists(rel, col string) Error {
	return Error{
		Code:    "42701",
		Message: fmt.Sprintf("column \"%s\" of relation \"%s\" already exists", col, rel),
	}
}

func ErrorColumnDoesNotExist(rel, col string) Error {
	return Error{
		Code:    "42703",
		Message: fmt.Sprintf("column \"%s\" of relation \"%s\" does not exist", col, rel),
	}
}

func ErrorRelationAlreadyExists(rel string) Error {
	return Error{
		Code:    "42P07",
		Message: fmt.Sprintf("relation \"%s\" already exists", rel),
	}
}

func ErrorRelationDoesNotExist(tbl string) Error {
	return Error{
		Code:    "42P01",
		Message: fmt.Sprintf("relation \"%s\" does not exist", tbl),
	}
}

func ErrorSchemaAlreadyExists(sch string) Error {
	return Error{
		Code:    "42P06",
		Message: fmt.Sprintf("schema \"%s\" already exists", sch),
	}
}

func ErrorSchemaDoesNotExist(sch string) Error {
	return Error{
		Code:    "3F000",
		Message: fmt.Sprintf("schema \"%s\" does not exist", sch),
	}
}

func ErrorTypeAlreadyExists(typ string) Error {
	return Error{
		Code:    "42710",
		Message: fmt.Sprintf("type \"%s\" already exists", typ),
	}
}

func ErrorTypeDoesNotExist(typ string) Error {
	return Error{
		Code:    "42704",
		Message: fmt.Sprintf("type \"%s\" does not exist", typ),
	}
}

// severity: ERROR
// code: 42701
// message: column "bar" of relation "foo" already exists
