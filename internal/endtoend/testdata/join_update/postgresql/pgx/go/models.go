// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package querytest

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type GroupCalcTotal struct {
	Npn     pgtype.Text
	GroupID pgtype.Text
}

type NpnExternalMap struct {
	ID  pgtype.Text
	Npn pgtype.Text
}

type ProducerGroupAttribute struct {
	NpnExternalMapID pgtype.Text
	GroupID          pgtype.Text
}
