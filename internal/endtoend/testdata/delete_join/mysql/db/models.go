// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package db

import ()

type JoinTable struct {
	ID             int64
	PrimaryTableID int64
	OtherTableID   int64
	IsActive       bool
}

type PrimaryTable struct {
	ID     int64
	UserID int64
}
