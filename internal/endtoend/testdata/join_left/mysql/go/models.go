// Code generated by sqlc. DO NOT EDIT.

package querytest

import (
	"database/sql"
)

type Author struct {
	ID       int32
	Name     string
	ParentID sql.NullInt32
}

type City struct {
	CityID  int32
	MayorID int32
}

type Mayor struct {
	MayorID  int32
	FullName string
}

type SuperAuthor struct {
	SuperID       int32
	SuperName     string
	SuperParentID sql.NullInt32
}

type User struct {
	UserID int32
	CityID sql.NullInt32
}
